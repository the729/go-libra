var defaultServer = "http://hk2.wutj.info:38080";
var client = libra.client(defaultServer, libra.trustedPeersFile);

const fromHexString = hexString =>
    new Uint8Array(hexString.match(/.{1,2}/g).map(byte => parseInt(byte, 16)));

const toHexString = bytes =>
    bytes ? bytes.reduce((str, byte) => str + byte.toString(16).padStart(2, '0'), '') : 'N/A';

Vue.component("ledger-info", {
    template: `
    <div class="ledger-info">
        <h2>Ledger Info</h2>
        <ul v-if="ledger_info">
            <li>Ledger Version: <span>{{ ledger_info.getVersion() }}</span></li>
            <li>Timestamp: <span>{{ ledger_info.getTimestampUsec() }} ({{ new Date(ledger_info.getTimestampUsec()/1000.0).toLocaleString() }})</span></li>
            <li>Epoch Number: {{ ledger_info.getEpochNum() }}</li>
            <li>Accumulator Hash: {{ toHexString(ledger_info.getTransactionAccumulatorHash()) }}</li>
        </ul>
        <button @click="refresh" :disabled="loading">{{ loading?"Refreshing...":"Refresh" }}</button>
    </div>
    `,
    data: function () {
        return {
            ledger_info: null,
            loading: false
        }
    },
    methods: {
        refresh: function () {
            this.loading = true;
            client.queryLedgerInfo()
                .then(r => {
                    this.ledger_info = r;
                    this.$emit('ledger-info-updated', r);
                })
                .catch(e => { console.log("Error:", e) })
                .finally(_ => { this.loading = false })
        }
    },
    mounted: function () {
        this.refresh()
    }
})

Vue.component("txn-card", {
    props: ["txn"],
    template: `
    <tr @click="$emit('txn-selected', txn)" class="pointer">
        <td>{{ version }}</td>
        <td>{{ type }}</td>
        <td>{{ exp_time }}<span v-if="isFinite(exp_time)"><br/>({{ new Date(exp_time*1000).toLocaleTimeString() }})</span></td>
        <td>Sender: 
            <a @click.stop="$emit('addr-selected', toHexString(sender))">{{ toHexString(sender) }}</a>
            <br/>Receiver: 
            <a @click.stop="$emit('addr-selected', toHexString(receiver))">{{ toHexString(receiver) }}</a>
        </td>
        <td>{{ amount }}</td>
        <td>{{ vm_result }}</td>
        <td>{{ gas_used }}</td>
    </tr>
    `,
    computed: {
        version: function () { return this.txn.getVersion() },
        gas_used: function () { return this.txn.getGasUsed() },
        vm_result: function () { return this.txn.getMajorStatus() },
        raw_txn: function () {
            var stxn = this.txn.getSignedTxn();
            return stxn ? stxn.RawTxn : null
        },
        block_meta: function () { return this.txn.getBlockMetadata() },
        exp_time: function () {
            return this.raw_txn ? this.raw_txn.ExpirationTime :
                this.block_meta ? this.block_meta.TimestampUSec / 1000000 : "N/A";
        },
        type: function () {
            if (this.block_meta) return "block_meta";
            if (!this.raw_txn) return "unknown";
            if (this.raw_txn.Payload.WriteSet) return "write_set";
            if (!this.raw_txn.Payload.Code) return "unknown";
            var name = libra.inferProgramName(this.raw_txn.Payload.Code);
            switch (name) {
                case "peer_to_peer_transfer": return "p2p";
                case "create_account": return "create account";
                case "mint": return "mint";
                case "rotate_authentication_key": return "rotate key";
            }
            return "unknown";
        },
        sender: function () {
            return this.raw_txn ? this.raw_txn.Sender :
                this.block_meta ? this.block_meta.Proposer : null;
        },
        receiver: function () {
            if (!this.raw_txn) return null;
            switch (this.type) {
                case "p2p":
                case "mint":
                    return this.raw_txn.Payload.Args[0];
            }
            return null;
        },
        amount: function () {
            if (!this.raw_txn) return null;
            switch (this.type) {
                case "p2p":
                case "mint":
                    return this.raw_txn.Payload.Args[1] / 1000000.0;
            }
            return null;
        }
    }
})

Vue.component("txn-list", {
    template: `
    <div class="txn-list">
        <h2>Transaction List</h2>
        <div class="controls">
            Range: <input v-model="version_from" :disabled="latest_10"/> - <input v-model="version_to" :disabled="latest_10"/>
            <input type="checkbox" id="latest_10_checkbox" v-model="latest_10">Latest 10 Txns</input>
            <button @click="refresh" :disabled="loading">{{ loading?"Refreshing...":"Refresh" }}</button>
        </div>
        <table>
            <thead>
                <tr>
                    <th>Version</th>
                    <th>Type</th>
                    <th>Exp. Time</th>
                    <th>Sender / Receiver</th>
                    <th>Amount</th>
                    <th>Result</th>
                    <th>Gas Used</th>
                </tr>
            </thead>
            <tbody>
                <txn-card v-for="txn in txns"
                    :key="txn.getVersion()"
                    :txn="txn"
                    @txn-selected="$emit('txn-selected', $event)"
                    @addr-selected="$emit('addr-selected', $event)"
                />
            </tbody>
        </table>
    </div>
    `,
    props: ['ledger_info'],
    data: function () {
        return {
            latest_10: true,
            version_from: 0,
            version_to: 0,
            txns: [],
            loading: false
        }
    },
    watch: {
        latest_10: function (val) {
            if (val) {
                this.version_from = this.ledger_info.getVersion();
                this.version_to = this.version_from - 9;
                if (this.version_to < 0) this.version_to = 0;
            }
        },
        ledger_info: function (val) {
            if (this.latest_10) {
                this.version_from = val.getVersion();
                this.version_to = this.version_from - 9;
                if (this.version_to < 0) this.version_to = 0;
                this.refresh();
            }
        }
    },
    methods: {
        refresh: function () {
            var start, limit, asc;
            if (this.version_from > this.version_to) {
                limit = this.version_from - this.version_to + 1;
                asc = false;
                if (limit > 25) {
                    limit = 25;
                    this.version_to = parseInt(this.version_from) - 24;
                }
                start = this.version_to;
            } else {
                start = this.version_from;
                limit = this.version_to - this.version_from + 1;
                asc = true;
                if (limit > 25) {
                    limit = 25;
                    this.version_to = parseInt(this.version_from) + 24;
                }
            }
            this.loading = true;
            client.queryTransactionRange(start, limit, true)
                .then(r => {
                    var txns = r.getTransactions()
                    if (!asc) {
                        txns.reverse()
                    }
                    this.txns = txns;
                })
                .catch(e => { console.log("Error:", e) })
                .finally(_ => { this.loading = false })
        }
    }
})

Vue.component("txn-detail", {
    props: ['txn'],
    template: `
    <div class="txn-detail">
        <h2>Transaction # {{ version }} Details</h2>
        <ul>
            <li>Gas used (micro Libra): {{ gas_used }}</li>
            <li>Major status: {{ vm_result }}</li>
            <li v-if="raw_txn != null"> Transaction content:
                <ul>
                    <li>Sender: 
                        <a @click="$emit('addr-selected', toHexString(raw_txn.Sender))">{{ toHexString(raw_txn.Sender) }}</a>
                    </li>
                    <li>Sender Seq. Number: {{ raw_txn.SequenceNumber }}</li>
                    <li>Max Gas Amount: {{ raw_txn.MaxGasAmount }} gas</li>
                    <li>Gas Price: {{ raw_txn.GasUnitPrice }} ulibra/gas</li>
                    <li>Expiration: {{ raw_txn.ExpirationTime }}
                    ({{ new Date(raw_txn.ExpirationTime * 1000).toLocaleString() }})</li>
                    <li>Payload:
                        <ul>
                            <li>Script Type: {{ type }}</li>
                            <li v-if="type=='mint'||type=='p2p'">Receiver: 
                                <a @click="$emit('addr-selected', toHexString(receiver))">{{ toHexString(receiver) }}</a>
                            </li>
                            <li v-if="type=='mint'||type=='p2p'">Amount: {{ amount }} libra</li>
                        </ul>
                    </li>
                </ul>
            </li>
            <li v-else-if="block_meta != null"> Block metadata:
                <ul>
                    <li>ID: {{ toHexString(block_meta.ID) }}</li>
                    <li>Proposer: {{ toHexString(block_meta.Proposer)}}</li>
                    <li>Timestamp usec: {{ block_meta.TimestampUSec }}</li>
                </ul>
            </li>
            <li v-else>Not a user transaction.</li>
        </ul>
        <p>(See console for raw SignedTxn object.)</p>
    </div>
    `,
    computed: {
        version: function () { return this.txn.getVersion() },
        gas_used: function () { return this.txn.getGasUsed() },
        vm_result: function () { return this.txn.getMajorStatus() },
        raw_txn: function () {
            var stxn = this.txn.getSignedTxn();
            return stxn ? stxn.RawTxn : null
        },
        block_meta: function () { return this.txn.getBlockMetadata() },
        type: function () {
            if (this.block_meta) return "block_meta";
            if (!this.raw_txn) return "unknown";
            if (this.raw_txn.Payload.WriteSet) return "write_set";
            if (!this.raw_txn.Payload.Code) return "unknown";
            var name = libra.inferProgramName(this.raw_txn.Payload.Code);
            switch (name) {
                case "peer_to_peer_transfer": return "p2p";
                case "create_account": return "create account";
                case "mint": return "mint";
                case "rotate_authentication_key": return "rotate key";
            }
            return "unknown";
        },
        receiver: function () {
            if (!this.raw_txn) return null;
            switch (this.type) {
                case "p2p":
                case "mint":
                    return this.raw_txn.Payload.Args[0];
            }
            return null;
        },
        amount: function () {
            if (!this.raw_txn) return null;
            switch (this.type) {
                case "p2p":
                case "mint":
                    return this.raw_txn.Payload.Args[1] / 1000000.0;
            }
            return null;
        }
    }
})

Vue.component("account-detail", {
    props: ['address'],
    template: `
    <div class="account-detail">
        <h2>Account Detail</h2>
        <p>
            Address: <input type="text" v-model="address"></input>
            <button @click="refresh" :disabled="loading">{{ loading?"Refreshing...":"Refresh" }}</button>
        </p>
        <div v-if="state">
            <div v-if="blob">
                <p>All resources:</p>
                <ul>
                    <li v-for="path in paths">{{ toHexString(path) }} (Type:{{ path[0]==0?"Code":path[0]==1?"Resource":"Unknown" }})<br/>
                        Raw content: {{ toHexString(blob.getResource(path)) }}
                    </li>
                </ul>
                <p>Libra account resource:</p>
                <ul v-if="resource">
                    <li>Balance: {{ resource.getBalance() / 1000000.0 }}</li>
                    <li>Sequence Number: {{ resource.getSequenceNumber() }}</li>
                    <li>Sent Events: 
                        <ul>
                            <li>Count: {{ resource.getSentEvents().Count }}</li>
                            <li>Key: {{ toHexString(resource.getSentEvents().Key) }}</li>
                        </ul>
                    </li>
                    <li>Received Events: 
                        <ul>
                            <li>Count: {{ resource.getReceivedEvents().Count }}</li>
                            <li>Key: {{ toHexString(resource.getReceivedEvents().Key) }}</li>
                        </ul>
                    </li>
                    <li>Delegated withdrawal capability: {{ resource.getDelegatedWithdrawalCapability() }}</li>
                    <li>Event generator: {{ resource.getEventGenerator() }}</li>
                    <li>Proven at ledger version: {{ resource.getLedgerInfo().getVersion() }}</li>
                </ul>
                <p v-else>Account exists but cannot get 0x0.LibraAccount.T resource.</p>
            </div>
            <p v-else>Account does not exist at version {{ state.getVersion() }}</p>
        </div>
    </div>
    `,
    data: function () {
        return {
            loading: false,
            state: null,
            blob: null,
            resource: null,
            paths: []
        }
    },
    methods: {
        refresh: function () {
            if (this.address.length != 64) {
                alert("Invalid address.");
                return
            }

            var addr = fromHexString(this.address);

            this.loading = true;
            client.queryAccountState(addr)
                .then(r => {
                    this.state = r;
                    return r.getAccountBlob()
                })
                .then(r => {
                    this.blob = r;
                    if (r == null) return null;
                    this.paths = r.getResourcePaths()
                    return r.getLibraAccountResource()
                })
                .then(r => {
                    this.resource = r;
                    setTimeout(_ => {
                        window.scrollTo(0, this.$el.offsetTop)
                    }, 10)
                })
                .catch(e => {
                    console.log("Error: ", e)
                })
                .finally(_ => { this.loading = false })
        }
    }
})

new Vue({
    el: "#app",
    data: {
        selected_txn: null,
        ledger_info: null,
        selected_address: null
    },
    methods: {
        li_update: function (li) {
            console.log("Ledger info", li)
            this.ledger_info = li;
        },
        txn_selected: function (txn) {
            this.selected_txn = txn;
            if (txn.getSignedTxn()) console.log("Selected Txn (user txn): ", txn.getSignedTxn())
            if (txn.getBlockMetadata()) console.log("Selected Txn (block meta): ", txn.getBlockMetadata())
            setTimeout(_ => window.scrollTo(0, this.$el.querySelector("#txn-detail").offsetTop), 10)
        },
        addr_selected: function (addr) {
            this.selected_address = addr;
            setTimeout(_ => {
                this.$refs.account_detail.refresh();
                window.scrollTo(0, this.$el.querySelector("#account-detail").offsetTop)
            }, 10)
        }
    }
})

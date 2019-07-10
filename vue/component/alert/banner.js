Vue.component('alert-banner', {
    props: ['alerts'],

    data: function () {
        return {
            showReload: false
        }
    },

    watch: { 
        alerts: function(newVal, oldVal) { 
            for (let i = 0; i < this.alerts.length; i++) { 
                if (this.alerts[i].error.message.indexOf("invalid session") == -1 ||
                    this.alerts[i].type === "connection") {
                    this.showReload = true;
                }
            }
        }
    },

    template: `
<div v-if="alerts.length > 0" class="box has-background-danger has-text-white content">
    <a v-if="showReload" href="" class="button is-pulled-right">Try Reload</a>
    <p class="has-text-white has-text-weight-bold is-size-4">Alert</p>

    <li v-for="a in alerts">{{ a.error.message }}</li>
</div>`
})
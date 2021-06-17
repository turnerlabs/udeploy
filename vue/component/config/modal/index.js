Vue.component('config-modal', {
    template: '#config-modal-template',
    
    props: ['app', 'instance', 'link'],

    data: function() {
        let data = {
            value: "",
           
            pending: false,

            saveAuthorized: (this.instance.deployCode.length == 0),

            error: "",
            warn: "",
        }

        return data
    },

    mounted: function () { 
        this.loadValue(this.link.linkId);
    },

    methods: {
        loadValue: function (linkId) {
            let that = this;

            fetch('/v1/apps/'+this.app.name+"/instances/"+this.instance.name+"/config?linkId="+linkId)
            .then(function(response) {
                return response.json()
            })
            .then(function(data) {
                if (data.message) {
                    return Promise.reject(new Error(data.message))
                }

                return Promise.resolve(data);
            })
            .then(function(data) {
                try {
                    that.value = JSON.stringify(JSON.parse(data.value), null, 4);
                } catch(e) {
                    that.value = data.value;
                }
            })
            .catch(function(e) {
                that.error = e.message;
            });
        },
        saveValue: function () {
            this.error = ""
            this.pending = true
           
            let that = this

            let body = {
                linkId: this.link.linkId,
                value: this.value
            }

            return fetch('/v1/apps/' + this.app.name + "/instances/" + this.instance.name + "/config", {
                method: "POST",
                body: JSON.stringify(body),
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json' 
                }
            })
            .then(function(response) {
                return response.json()
            })
            .then(function(data) { 
                if (data && data.message) {
                    that.error = data.message;
                } else {   
                    that.$parent.modal.config.show = false 
                }

                that.pending = false
            });
        },
        validateDeployCode: function (evt) {
            this.saveAuthorized = (this.instance.deployCode == evt.target.value)
        }
    }
})
Vue.component('scale-modal', {
    template: '#scale-modal-template',
    
    props: ['instance', 'app', 'action'],

    data: function() {
        let data = {
            count: this.instance.desiredCount,
            code: "",

            startAuthorized: false,

            starting: false,

            error: ""
        }

        return data
    },

    mounted: function () {
        this.validateStart()
    },

    watch: {
        count: function (val) {
            this.validateStart()
        },
    },

    methods: {
        validateStart: function (evt) {
            this.instance.deployCode.length > 0
                ? this.startAuthorized = this.validateCount() && this.validateCode()
                : this.startAuthorized = this.validateCount()
        },
        validateCode: function () {
            return this.instance.deployCode == this.code
        },
        validateCount: function () {
            return this.count > 0 && this.count <= 50
        },
        scale: function() {
            let that = this
            
            this.starting = true

            let body = {
                "version": this.instance.version
            }

            fetch('/v1/apps/' + this.app.name + '/instances/' + this.instance.name + '/scale/' + this.count, {
                method: "PUT",
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
                if (data.message) {
                    that.error = data.message
                    that.starting = false
                } else {
                    that.instance.deployment.isPending = true;

                    that.$parent.updateInstance(that.app.name, that.instance.name, that.instance)

                    that.$parent.modal.scale.show = false
                } 
            })
        }
    }
})
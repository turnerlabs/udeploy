Vue.component('confirm-modal', {
    template: '#confirm-modal-template',
    
    props: ['message', 'code'],

    data: function() {
        let data = {
            confirming: false,
            
            authorized: (this.code.length == 0)
        }

        return data
    },

    methods: {
        validateCode: function (evt) {
            this.authorized = (this.code == evt.target.value)
        },
        confirm: function() {
            this.$emit('continue');
            this.$parent.modal.confirm.show = false
        }
    }
})
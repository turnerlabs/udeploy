export let url = {
    current: function (url) {
        return {
            params: (template) => {

                let tokens = template.split("/")
                let values = url.split("/")
            
                let uri = {}
            
                tokens.map((t, i) => {
                    if (t.indexOf(":") == 0) {
                        name = t.trimLeft(":")
            
                        uri[t.substring(1,name.length)] = values[i]
                    }
                })
            
                return uri
            }
        }
    }(window.location.pathname)
}
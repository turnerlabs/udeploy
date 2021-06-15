import {url} from "../../component/url/params.js";
import {obj} from "../../component/copy/object.js";
import {includeTenplates} from "../../component/html/include.js";

includeTenplates().then(() => {
    var app = new Vue({
        el: '#app',

        data: {
            params: url.current.params("/apps/:app"),
            
            user: {email:""},
            config: {consoleLink:""},

            isLoading: false,
            isDeleting: false,
            
            canSave: false,
            alerts: [],

            app: [],
            projects: [],

            modal: {
                confirm:  {
                    show: false,
                    message: "",
                    code: "",
                    action: ()=>{},
                },
                image:  {
                    show: false,
                    message: "",
                    img: "",
                }
            },

            defaultApp: {
                "name": "application-name",
                "type": "service",
                "project": { "id": "000000000000000000000000"},
                "repo": { 
                    "tagFormat": "version",
                    "commitConfig" : {
                        "limit": 20,
                    },
                },
                "instances": [{
                    "edited": true,
                    "name":"dev",
                    "order":1,
                    "cluster":"",
                    "service":"",
                    "repository":"",
                    "deployCode":"",
                    "links": [],
                    "task":{
                        "family":"",
                        "registry":"",
                        "imageTagEx":"",
                        "revisions":50,
                    },
                    "claims": {
                        "edit": true,
                    }  
                }],
            },  

            lambdaRegex: "(v[0-9]+.[0-9]+.[0-9]+-[\\w]+)[.]?([a-zA-Z0-9]+)?",
            serviceRegex: "(v[0-9]+.[0-9]+.[0-9]+-[a-zA-Z0-9]*).([0-9]+)",

            page: {}
        },

        created: function () {
            let that = this;
            
            that.isLoading = true;

            this.getConfig();

            this.getProjects().then(function(projects) {
                that.projects = projects;
            });

            this.getUser().then(function(user) {
                that.user = user

                if (!that.isNew()) {
                    
                    fetch('/v1/apps/'+that.params.app+"?configOnly=true").then(function(response) {
                        return response.json()
                    })
                    .then(function(data) {
                        if (data.message) {
                            return Promise.reject(new Error(data.message))
                        }
    
                        return Promise.resolve(data);
                    })
                    .then(function(app) {

                        app.instances = app.instances.sort((a, b) => (a.order > b.order) ? 1 : -1)

                        for (let x = 0; x < app.instances.length; x++) {
                            app.instances[x].links = app.instances[x].links.sort((a, b) => (a.generated < b.generated) ? 1 : -1)
                            app.instances[x].links = that.addUniqueIds(app.instances[x].links)
                        }
                        
                        that.app = app;
                       
                        that.canSave = that.userCanSave();

                        that.isLoading = false;
                    })
                    .catch(function(e) {
                        that.addError(e.message, "general");
                    }); 
                } else {
                    that.app = obj.copy(that.defaultApp);
                    that.isLoading = false;
                }
            })
            .catch(function(e) {
                that.addError(e.message, "general");
            });    
        },

        methods: {
            isNew: function() {
                return (this.params.app.toLowerCase() == "new")
            },
            handleContinue: function () {
                this.modal.confirm.action()
            },
            removeLink(index, linkIndex) {
                this.app.instances[index].links.splice(linkIndex, 1);
            },
            addLink(index) {
                this.app.instances[index].links.push(
                {
                    "id":  this.newUniqueId(this.app.instances[index].links)
                });
            },
            addUniqueIds(list) {
                for (let y = 0; y < list.length; y++) {
                    list[y].id = y+1;
                }
    
                return list
            },
            newUniqueId(list) {
                let max = 1;
    
                list.forEach(item => {
                    if (item.id >= max) {
                        max = item.id + 1
                    }
                });

                return max
            },
            validateField(evt, regex, msg, type) {
                var re = new RegExp(regex);
                
                if (re.test(evt.target.value)) {
                    evt.target.classList.remove("is-danger")
                    this.removeError(msg, type)

                } else {
                    evt.target.classList.add("is-danger")
                    this.addError(msg, type)
                }

                this.canSave = this.userCanSave();
            },
            removeErrorsByType(type) {
                for( var i = 0; i < this.alerts.length; i++){ 
                    if (this.alerts[i].type === type) {
                      this.alerts.splice(i, 1); 
                      i--
                    }
                 }
            },
            removeError(msg, type) {
                for( var i = 0; i < this.alerts.length; i++){ 
                    if (this.alerts[i].error.message === msg && this.alerts[i].type === type) {
                      this.alerts.splice(i, 1); 
                      i--
                    }
                 }
            },
            addError(msg, type) {
                if (this.notReported(msg, type)) {
                    this.alerts.push({
                        error: new Error(msg),
                        type: type
                    })
                }
            },
            notReported(msg, type) {
                for( var i = 0; i < this.alerts.length; i++){ 
                    if (this.alerts[i].error.message === msg && this.alerts[i].type === type ) {
                      return false;
                    }
                 }

                 return true
            },
            instanceNameChanged(name) {
                for (let x = 0; x < this.app.instances.length; x++) {
                    if (this.app.instances[x].name === name) {
                        this.app.instances[x].edited = true;
                    }
                } 
            },
            userCanSave: function() {
                for (let x=0; x < this.alerts.length; x++) {
                    if (this.alerts[x].type != "general") {
                        return false
                    }
                }

                if (this.user.admin) {
                    return true
                }

                if (!this.user.apps[this.app.name].claims || this.user.apps[this.app.name].claims.length == 0) {
                    return false
                }

                for (const instance of Object.keys(this.user.apps[this.app.name].claims)) {
                    for (const claim of this.user.apps[this.app.name].claims[instance]) {
                        if (claim === "edit") {
                            return true
                        }
                    }
                }

                return false
            },
            deleteInstance: function(index) {
                this.removeErrorsByType(index) 
                this.canSave = this.userCanSave();

                Vue.delete(this.app.instances, index)
            },
            addInstance: function() {
                let d = obj.copy(this.defaultApp.instances[0]);

                for (let x = 0; x < this.app.instances.length; x++) {
                    if (this.app.instances[x].claims["edit"] || this.user.admin) {
                        d = obj.copy(this.app.instances[x]);
                        d.name = d.name + "-"+ this.app.instances.length
                        d.task.tasksInfo = [];
                        d.links = [];
                        d.edited = true;
                    }
                }

                Vue.set(this.app.instances, this.app.instances.length, d)
            },
            cancel: function() {
                window.location.href = "/apps";
            },
            save: function() {
                let that = this;

                this.removeErrorsByType("general")
                
                if (this.alerts.length > 0) {
                    return
                }

                return fetch('/v1/apps/' + this.params.app, {
                    method: "POST",
                    body: JSON.stringify(this.app),
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
                        return Promise.reject(new Error(data.message))
                    }

                    return Promise.resolve(data);
                })
                .then(function(data) { 
                    window.location.href = "/apps#"+that.app.name;
                })
                .catch(function(e){
                    that.addError(e.message, "general");
                });
            },
            trash() {
                let that = this
                
                that.isDeleting = true

                fetch('/v1/apps/' + this.app.id, { method: "DELETE" })
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
                    that.isDeleting = false

                    window.location.href = "/apps";
                })
                .catch(function(e){
                    that.addError(e.message, "general");
                });

            },
            sortInstances: function (instances) {
                let temp = obj.copy(instances)
                
                return temp.sort(function(a,b) {
                    return a.order - b.order;
                });
            },
            getUser: function() {
                let that = this

                return fetch('/v1/user')
                .then(function(response) {
                    return response.json()
                })
                .then(function(data) {
                    if (data.message) {
                        return Promise.reject(new Error(data.message))
                    }

                    return Promise.resolve(data);
                })
            },
            getProjects: function() {
                let that = this

                return fetch('/v1/projects')
                .then(function(response) {
                    return response.json()
                })
                .then(function(data) {
                    if (data.message) {
                        return Promise.reject(new Error(data.message))
                    }

                    return Promise.resolve(data);
                })
            },
            getConfig: function() {
                let that = this

                fetch('/config')
                .then(function(response) {
                    return response.json()
                })
                .then(function(data) {
                    if (data.message) {
                        return Promise.reject(new Error(data.message))
                    }

                    that.config = data
                }).catch(function(e) {
                    that.alerts.push({ error: e });
                });    
            },
        }
    })
})

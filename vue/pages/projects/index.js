import {obj} from "../../component/copy/object.js";
import {includeTenplates} from "../../component/html/include.js";

includeTenplates().then(() => {
  var projects = new Vue({
    el: '#projects',

    data: {
        user: {email:""},
        config: {consoleLink:""},

        projects: [],
        apps: [],

        isLoading: false,
        canSave: false,
        new: true,
        alerts: [],

        deletedProjects:[],

        defaultProject: {
            "name": "New Project",
            "apps": []
        },

        page: {
            projects: true,
        }
    },

    created: function () {
        let that = this;

        that.isLoading = true;

        this.getConfig();

        this.getUser().then(function(user) {
            that.user = user;
            
            fetch('/v1/apps?all=true').then(function(response) {
                return response.json()
            })
            .then(function(data) {
                if (data.message) {
                    return Promise.reject(new Error(data.message))
                }
    
                return Promise.resolve(data);
            })
            .then(function(apps) {
                that.apps = apps;

                fetch('/v1/projects').then(function(response) {
                    return response.json()
                })
                .then(function(data) {
                    if (data.message) {
                        return Promise.reject(new Error(data.message))
                    }
    
                    return Promise.resolve(data);
                })
                .then(function(projects) {
                    
                    for (let x=0; x < projects.length; x++) {
                        projects[x].apps = [];

                        for (let y=0; y < that.apps.length; y++) {
                            if (that.apps[y].project.id === projects[x].id) {
                                projects[x].apps.push(that.apps[y])
                            }
                        }
                    }

                    that.projects = projects;

                    that.updateSaveState();
                   
                    that.isLoading = false;
                }).catch(function(e) {
                    that.addError(e.message, "general");
    
                    that.isLoading = false;
                });  
            })
            .catch(function(e) {
                that.addError(e.message, "general");
            });  
        }).catch(function(e) {
            that.addError(e.message, "general");
            
            that.isLoading = false;
        });  
       
       
    },

    methods: {
        validateField(evt, regex, msg, type) {
            var re = new RegExp(regex);

            if (re.test(evt.target.value)) {
                evt.target.classList.remove("is-danger")
                this.removeError(msg, type)
            } else {
                evt.target.classList.add("is-danger")
                this.addError(msg, type)
            }

            this.updateSaveState();
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
                    type: type,
                })
            }
        },
        notReported(msg, type) {
            for( var i = 0; i < this.alerts.length; i++){ 
                if (this.alerts[i].error.message === msg && this.alerts[i].type === type) {
                  return false;
                }
             }

             return true
        },
        didEdit: function(idx) {
            this.projects[idx].edited = true;
        },
        updateSaveState: function() {
            this.validate();

            if (this.alerts.length > 0) {
                this.canSave = false
            } else {
                this.canSave = this.user.admin;
            }
        },
        validate: function() {
            this.removeErrorsByType("general")
        },
        save: function() {
            let that = this;

            if (this.alerts.length > 0) {
                return
            }

            let actions = [];

            for (let x=0; x < this.deletedProjects.length; x++) {
                actions.push(fetch('/v1/projects/' + this.deletedProjects[x].id, {
                    method: "DELETE"
                }))
            }
    
            for (let x=0; x < this.projects.length; x++) {
                if (!this.projects[x].edited) {
                    continue
                }
                
                actions.push(fetch('/v1/projects/' + this.projects[x].id, {
                    method: "POST",
                    body: JSON.stringify(this.projects[x]),
                    headers: {
                        'Accept': 'application/json',
                        'Content-Type': 'application/json' 
                    }
                }))
            }

            Promise.all(actions)
            .then(responses => Promise.all(
                responses.map(r => r.json())
            ))
            .then(function(responses) {
                for (let x=0; x< responses.length; x++) {
                    if (responses[x].message) {
                        that.addError(responses[x].message, "general");
                    }
                }

                if (that.alerts.length == 0) {
                    window.location.href = "/apps";
                }
            })
            .catch((err) => {  
                that.AddError(err, "general");
            });
        },
        cancel: function() {
            window.location.href = "/apps";
        },
        deleteProject: function(index) {
            this.removeErrorsByType(index)
            
            this.updateSaveState()

            if (this.projects[index].id) {
                this.deletedProjects.push(this.projects[index])
            }
            
            Vue.delete(this.projects, index)
        },
        addProject: function() {
            let newProject = obj.copy(this.defaultProject)

            newProject.edited = true;
            
            Vue.set(this.projects, this.projects.length, newProject)
        },
        getUser: function() {
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


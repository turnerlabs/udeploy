import {obj} from "../../component/copy/object.js";
import {includeTenplates} from "../../component/html/include.js";

includeTenplates().then(() => {
  var users = new Vue({
    el: '#users',

    data: {
        user: {email:""},
        config: {consoleLink:""},

        apps: [],

        isLoading: false,
        canSave: false,
        new: true,
        alerts: [],

        users: [],
        deletedUsers:[],

        defaultUser: {
            "email": "user@domain.com",
            "apps": {
                "application-name": {
                    "claims": {
                        "dev": [
                            "scale",
                            "deploy",
                            "edit"
                        ]
                    }
                }
            }
        },

        page: {
            users: true,
        }
    },

    created: function () {
        let that = this;

        that.isLoading = true;

        this.getConfig();

        this.getUser().then(function(user) {
            that.user = user
            
            fetch('/v1/users').then(function(response) {
                return response.json()
            })
            .then(function(data) {
                if (data.message) {
                    return Promise.reject(new Error(data.message))
                }

                return Promise.resolve(data);
            })
            .then(function(users) {
                for (let x=0; x < users.length; x++) {
                    users[x].policy = JSON.stringify(users[x].apps, null, 5);
                    users[x].showPolicy = false;
                }

                that.users = users;
                that.updateSaveState();
               
                that.isLoading = false;
            })
            .catch(function(e) {
                that.addError(e.message, "general")
    
                that.isLoading = false;
           });
       }).catch(function(e) {
            that.addError(e.message, "general")

            that.isLoading = false;
       });
       
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
        })
        .catch(function(e) {
            that.addError(e.message, "general")
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
        validatePolicy(evt, type) {
            let msg = "REQUIRED: User Policy - Must be a properly formatted json object.";

            try {
                JSON.parse(evt.target.value)
                evt.target.classList.remove("is-danger")
                this.removeError(msg, type)
            } catch (e) {
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
        didEdit: function(id) {
            for (let x = 0; x < this.users.length; x++) {
                if (this.users[x].id === id) {
                    this.users[x].edited = true;
                }
            } 
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
            
            let hasAdmin, isDup = false
            for (let x=0; x< this.users.length; x++) {
                if (this.users[x].admin) {
                    hasAdmin = true;
                }

                for (let y=0; y< this.users.length; y++) {
                    if (this.users[x].id != this.users[y].id && this.users[x].email === this.users[y].email) {
                        isDup = true;
                    }
                }
            }

            if (!hasAdmin) {
                this.addError("At least one user must be an admin.", "general")
            }

            if (isDup) {
                this.AddError("User email addresses must be unique.", "general")
            }
        },
        save: function() {
            let that = this;

            if (this.alerts.length > 0) {
                return
            }

            let actions = [];

            for (let x=0; x < this.deletedUsers.length; x++) {
                actions.push(fetch('/v1/users/' + this.deletedUsers[x].id, {
                    method: "DELETE"
                }))
            }
    
            for (let x=0; x < this.users.length; x++) {
                if (!this.users[x].edited) {
                    continue
                }
                
                this.users[x].apps = JSON.parse(this.users[x].policy)

                actions.push(fetch('/v1/users/' + this.users[x].id, {
                    method: "POST",
                    body: JSON.stringify(this.users[x]),
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
                that.AddError(err.message, "general");
            });
        },
        cancel: function() {
            window.location.href = "/apps";
        },
        deleteUser: function(index) {
            this.removeErrorsByType(index)
            this.updateSaveState()

            this.deletedUsers.push(this.users[index])
            Vue.delete(this.users, index)
        },
        addUser: function() {
            let newUser = obj.copy(this.defaultUser)

            newUser.edited = true;
            newUser.showPolicy = true;
            newUser.policy = JSON.stringify(this.defaultUser.apps, null, 5);

            Vue.set(this.users, this.users.length, newUser)
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


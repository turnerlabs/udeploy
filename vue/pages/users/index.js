import {obj} from "../../component/copy/object.js";
import {includeTenplates} from "../../component/html/include.js";

includeTenplates().then(() => {
  var vm = new Vue({
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
            "apps": [],
        },

        page: {
            users: true,
        }
    },

    created: function () {
        let that = this;

        that.isLoading = true;

        this.getConfig();

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

            that.getUser().then(function(user) {
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
                    users = that.sortUsers(users)

                    for (let x=0; x < users.length; x++) {   
                        users[x] = that.addPolicy(users[x])
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
        })
        .catch(function(e) {
            that.addError(e.message, "general")
       });
       
    },

    methods: {
        addPolicy: function (user) {
            user.policy = []
            user.selectedRole = "";
            user.errorMessage = "";
           
            for (let a=0; a < this.apps.length; a++) {
                let appName = this.apps[a].name;

                let app = {
                    name: appName,
                    instances: [],
                }

                let usrApp = user.apps[appName]
            
                app.view = usrApp ? true : false
                
                for (var c in this.apps[a].instances) {
                    let instName = this.apps[a].instances[c].name

                    let inst = {
                        order: this.apps[a].instances[c].order,
                        name: instName,
                    }

                    if (usrApp && usrApp["claims"]) {
                        if (usrApp.claims.hasOwnProperty(instName)) {
                            for (let i=0; i < usrApp.claims[instName].length; i++) {
                                inst[usrApp.claims[instName][i]] = true;
                            }  
                        }
                    }

                    app.instances.push(inst)
                }
               
                user.policy.push(app)
            }

            return user
        },
        sortUsers: function (users) { 
            return users.slice().sort(function(a,b) {
                if (a.email < b.email)
                    return -1;
                if ( a.email > b.email)
                    return 1;
                return 0;
            });
        },
        sortApps: function (apps) { 
            return apps.slice().sort(function(a,b) {
                if (a.name < b.name)
                    return -1;
                if ( a.name > b.name)
                    return 1;
                return 0;
            });
        },
        sortInstances: function (instances) { 
            return instances.slice().sort(function(a,b) {
                if (a.order < b.order)
                    return -1;
                if (a.order > b.order)
                    return 1;
                return 0;
            });
        },
        formatRole(usr) {
            return (usr.roles && usr.roles.length > 0)
                ? usr.roles[0]
                : ""
        },
        getRole(usr) {
            return (usr.roles && usr.roles.length == 1) 
            ? usr.roles[0]
            : ""
        },
        inheritPolicy(usr) {
            usr.errorMessage = "";

            if (usr.selectedRole == "") {
                usr.roles = [];
                this.didEdit(usr.id);
                return
            }

            let previousRole = this.getRole(usr);

            if (this.isRoleValid(usr, {}, usr.selectedRole) && previousRole != usr.selectedRole) {
                usr.roles = [usr.selectedRole];

                this.didEdit(usr.id);
            } else {
                usr.errorMessage = "Inheriting the " + usr.selectedRole + " user policy would create an invalid circlular reference.";
            }
        },
        replicatePolicy(usr) {
            let u = this.findUser(usr.selectedRole, this.users);

            usr.policy = u.policy;

            this.didEdit(usr.id);
        },
        isRoleValid(usr, roles, testRole) {
            
            if (roles.hasOwnProperty(usr.email)) {
                return false
            }

            roles[usr.email] = true;

            if (testRole) {
                return this.isRoleValid(this.findUser(testRole, this.users), roles)
            } else if (usr.roles && usr.roles.length > 0) {  
                return this.isRoleValid(this.findUser(usr.roles[0], this.users), roles)
            }

            return true
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
                this.addError("User email addresses must be unique.", "general")
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

                this.users[x].apps = that.convertPolicy(this.users[x].policy)

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
        isInherited: function(email) {
            for (let x = 0; x < this.users.length; x++) {
                let role = this.getRole(this.users[x])

                if (role.length > 0 && role == email) {
                    return true
                }
            }

            return false
        }, 
        convertPolicy: function(policy) {
            let apps = {};
            
            for (let x=0; x< policy.length; x++) {
                let a = policy[x]

                if (!a.view) {
                    continue;
                }

                let app = { claims: {}};

                for (let i=0; i< a.instances.length; i++) {
                    app.claims[a.instances[i].name] = []

                    for (var c in a.instances[i]) {
                        if (c != "name" && c != "order" && a.instances[i][c]) {
                            app.claims[a.instances[i].name].push(c)
                        }
                    }
                }

                apps[policy[x].name] = app
            }

            return apps;
        },
        listUserApps: function (usr) {

            let apps = []

            if (usr.roles && usr.roles.length > 0) {
                let u = this.findUser(usr.roles[0], this.users);

                if (u != null) {
                    apps = this.listUserApps(u)
                }
            }
           
            return (usr.policy) 
                ? { ...(usr.policy
                    .filter((a) => a.view)
                    .reduce(function(map, obj) {
                        map[obj.name] = true;
                        return map;
                    }, {})), ...apps }
                : apps
        },
        findUser: function(email, users) {
            for (let x = 0; x < users.length; x++) {
                if (users[x].email == email) {
                    return users[x]
                }
            }

            return null
        },
        isUserMissing: function(email) {
            if (email == "") {
                return false
            }

            for (let x = 0; x < this.users.length; x++) {
                if (this.users[x].email == email) {
                    return false
                }
            }

            return true
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
            let newUser = this.addPolicy(obj.copy(this.defaultUser))

            newUser.edited = true;
            newUser.showPolicy = true;
            
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


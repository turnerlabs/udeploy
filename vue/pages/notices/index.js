import {obj} from "../../component/copy/object.js";
import {includeTenplates} from "../../component/html/include.js";

includeTenplates().then(() => {
  var notices = new Vue({
    el: '#notices',

    data: {
        user: {email:""},
        config: {consoleLink:""},

        notices: [],
        apps: [],
        instances: [],

        isLoading: false,
        canSave: false,
        new: true,
        alerts: [],

        deletedNotices:[],

        defaultNotice: {
            "name": "New Notice",
            "enabled": true,
            "showCriteria": true,
            "snsArn": "arn:aws:sns:us-east-1:{ACCOUNT}:{TOPIC}",
            "apps": [],
            "instances": [],
            "events": {
                "error": false,
                "starting": false,
                "pending": false,
                "running": false,
                "stopped": false,
                "deployed": false,
                "deploying": false,
            }
        },

        page: {
            notices: true,
        }
    },

    created: function () {
        let that = this;

        that.isLoading = true;

        this.getConfig();

        this.getUser().then(function(user) {
            that.user = user
            
            fetch('/v1/notices').then(function(response) {
                return response.json()
            })
            .then(function(data) {
                if (data.message) {
                    return Promise.reject(new Error(data.message))
                }

                return Promise.resolve(data);
            })
            .then(function(notices) {
                for (let x = 0; x < notices.length; x++) {
                    notices[x].showCriteria = false;

                    notices[x].instances = that.addUniqueIds(notices[x].instances);
                    notices[x].apps = that.addUniqueIds(notices[x].apps);
                } 

                that.notices = notices;
                that.updateSaveState();
               
                that.isLoading = false;
            }).catch(function(e) {
                that.addError(e.message, "general");

                that.isLoading = false;
            });  
        }).catch(function(e) {
            that.addError(e.message, "general");
            
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
                
            for (let x = 0; x < apps.length; x++) {
                for (let y = 0; y < apps[x].instances.length; y++) {
                    if (!that.instances.includes(that.apps[x].instances[y].name)) {
                        that.instances.push(that.apps[x].instances[y].name)
                    }
                }
            }
        })
        .catch(function(e) {
            that.addError(e.message, "general");
        });  
    },

    methods: {
        removeApp(type, index, appIndex) {
            this.notices[index].apps.splice(appIndex, 1);
            this.notices[index].edited = true;
            this.removeErrorsByType(type)
        },
        addApp(index) {
            this.notices[index].apps.push(
                {
                    "id":  this.newUniqueId(this.notices[index].apps),
                    "name": this.apps[0].name
                });
            this.notices[index].edited = true;

        },
        removeInst(type, index, instIndex) {
            this.notices[index].instances.splice(instIndex, 1);
            this.notices[index].edited = true;
            this.removeErrorsByType(type)
        },
        addInst(index) {
            this.notices[index].instances.push(
                {
                    "id":  this.newUniqueId(this.notices[index].instances),
                    "name": this.instances[0]
                });
            this.notices[index].edited = true;
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
            console.log(max)
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
            this.notices[idx].edited = true;
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

            for (let x=0; x < this.deletedNotices.length; x++) {
                actions.push(fetch('/v1/notices/' + this.deletedNotices[x].id, {
                    method: "DELETE"
                }))
            }
    
            for (let x=0; x < this.notices.length; x++) {
                if (!this.notices[x].edited) {
                    continue
                }
                
                actions.push(fetch('/v1/notices/' + this.notices[x].id, {
                    method: "POST",
                    body: JSON.stringify(this.notices[x]),
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
        deleteNotice: function(index) {
            this.removeErrorsByType(index)
            
            this.updateSaveState()

            if (this.notices[index].id) {
                this.deletedNotices.push(this.notices[index])
            }
            
            Vue.delete(this.notices, index)
        },
        addNotice: function() {
            let newNotice = obj.copy(this.defaultNotice)

            newNotice.edited = true;
            
            Vue.set(this.notices, this.notices.length, newNotice)
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


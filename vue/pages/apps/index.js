import {obj} from "../../component/copy/object.js";
import {includeTenplates} from "../../component/html/include.js";

includeTenplates().then(() => {
    var app = new Vue({
        el: '#app',

        data: {
            user: { email: "", admin: false },
            config: { consoleLink: "" },

            apps: [],
            modal: {
                deploy: {show: false},
                scale: {show: false, action: ""},
                confirm:  {
                    show: false,
                    message: "",
                    code: "",
                    action: ()=>{},
                }
            },

            selected: {
                app: {},
                instance: {},
                source: undefined,
            },

            dragged: {
                app: undefined,
                instance: undefined,
                targets: {},
            },

            filter: "",
            matches: {}, 

            version: "",
            updated: "",

            isLoading: true,
            
            alerts: [],

            page: {
                apps: true,
            },

            view: {
                grid: true,
                list: false,
            }
        },

        watch: {
            filter: function (val) {
                localStorage.setItem("filter", val);
            },
        },

        created: function () {
            this.getVersion()

            this.getConfig()

            let view = localStorage.getItem("view")

            if (view) {
                this.setView(view)
            }

            let filterValue = localStorage.getItem("filter")

            if (filterValue) {
                this.filter = filterValue
            }

            let that = this;
            this.getUser().then(function(user) {
                that.user = user

                that.refreshApps()

                that.watchForUpdates()
            }).catch(function(e) {
                that.alerts.push({ error: e });
            }); 
        },

        updated: function() {
           
            let hash = window.location.hash.replace("#", "")
            
            var elmnt = document.getElementById(hash);

            if (elmnt) {
                elmnt.scrollIntoView();
            }
        },

        methods: {
            formatErrorPreview(error) {
                if (error.length > 15) {
                    return error.substring(0, 20) + "..."
                }

                return error
            },
            setView(view) {
                this.view.list = false;
                this.view.grid = false;

                localStorage.setItem("view", view);
                this.view[view] = true;
            },
            formatVersion(version) {
                if (version.length > 16) {
                    return version.substring(0, 13) + "...";
                }

                return version
            },
            showApp(name) {
                if (this.filter.length == 0) {
                    Vue.set(this.matches, name, {})
                    return true
                }

                let index = name.indexOf(this.filter.toLowerCase())
               
                if (index != -1) {
                    Vue.set(this.matches, name, {})
                    return true
                }

                Vue.delete(this.matches, name)

                return false
            },
            hasEditPermission(app) {
                if (!this.user) {
                    return false
                }

                if (this.user.admin) {
                    return true
                }

                for (let i in app.instances) {
                    if (app.instances[i].claims.edit) {
                        return true
                    }
                }

                return false
            },
            addApp() {
                window.location.href = "/apps/new";
            },
            cacheApp(app) {
                let that = this
                
                that.apps.forEach(function (a, i) {
                    if (a.name == app.name) {
                        a.isRefreshing = true
                        Vue.set(that.apps, i, a)
                    }
                });

                fetch('/v1/apps/' + app.name + '/cache', { method: "PUT" })
                .then(function(response) {
                    return response.json()
                })
                .then(function(data) {
                    if (data.message) {
                        return Promise.reject(new Error(data.message))
                    }
                }).catch(function(e) {
                    that.alerts.push({ error: e });
                });
            },
            notDropTarget(app, inst) {
                return !inst.claims.hasOwnProperty("deploy") || (app.name != this.dragged.app.name || inst.name === this.dragged.instance.name)
            },
            dragStart(evt, app, item) {
                this.dragged.app = app;
                this.dragged.instance = item;

                app.instances.forEach((v)=>{
                    this.dragged.targets[v.name] = { counter: 0 };
                });
            },
            dragOver(evt, item) {},
            dragEnd(evt, item) {
                this.dragged.targets = {}
                this.dragged.app = undefined
                this.dragged.instance = undefined
            },
            dragEnter(evt, app, item) {
                if (this.notDropTarget(app, item)) {
                    return
                }

                this.dragged.targets[item.name].counter++
                
                if (evt.target.id == "box") {
                    
                    this.dragged.targets[item.name].el = evt.target;
                    this.dragged.targets[item.name].el.classList.add("can-drop");

                    this.selected.instance = item;
                    this.selected.app = app;
                }
            },
            dragLeave(evt, app, item) {
                if (this.notDropTarget(app, item)) {
                    return
                }

                this.dragged.targets[item.name].counter--
                if (this.dragged.targets[item.name].counter > 0) {
                    return
                }

                if (this.dragged.targets[item.name].el) {
                    this.dragged.targets[item.name].el.classList.remove("can-drop");
                }
                
                if (evt.target.id == "box") {
                    this.selected.instance = undefined;
                    this.selected.app = undefined;
                }
            },
            dragDrop(evt, app, item) {
                if (this.notDropTarget(app, item)) {
                    return
                }

                this.selected.source = this.dragged.instance.name;
                this.modal.deploy = { show: true };

                let el = evt.target
                while (true) {
                    if (el.id === "box") {
                        el.classList.remove("can-drop");
                        break;
                    }

                    el = el.parentNode;
                }
            },
            addError(err, type) {
                let add = true;

                for (let a in this.alerts) {
                    if (this.alerts[a].type == type || this.alerts[a].error.message == err.message) {
                        add = false;
                    } 
                } 

                if (add) {
                    this.alerts.push({ type: type, error: err });
                }
            },
            clearErrorsBy(type) {
                for (let a in this.alerts) {
                    if (this.alerts[a].type == type) {
                        this.alerts.splice(a, 1);
                    } 
                } 
            },
            watchForUpdates() {
                let that = this;

                if (!!window.EventSource) {
                    var source = new EventSource('events/app/changes');

                    source.onopen = function(){
                        console.log('connected');

                        that.clearErrorsBy("connection")
                    };

                    source.onerror = function(e) {
                        if (e.readyState == EventSource.CLOSED) {
                            console.log('closed')
                        }

                        that.addError(new Error("real-time updates connection lost"), "connection");
                    };
                    
                    source.addEventListener('message', function(e) {
                        let app = JSON.parse(e.data);

                        for (let a in that.apps) {
                            if (that.apps[a].name == app.name) {
                                Vue.set(that.apps, a, app)
                                that.updateTime();
                            }
                        }
                    }, false);
                } else {
                    that.alerts.push({ error: new Error("browser does not support real-time updates")})
                }
            },
            getInstance(appName, instanceName) {
                for (let a in this.apps) {
                    if (this.apps[a].name == appName) {
                        for (let i in this.apps[a].instances) {
                            if (this.apps[a].instances[i].name == instanceName) {
                                return this.apps[a].instances[i];
                            }
                        }
                    }
                }

                return undefined
            },
            updateInstance(appName, instanceName, instance) {
                for (let a in this.apps) {
                    if (this.apps[a].name == appName) {
                        for (let i in this.apps[a].instances) {
                            if (this.apps[a].instances[i].name == instanceName) {
                                this.apps[a].instances[i] = instance;
                            }
                        }
                    }
                }
            },
            updateTime: function () {
                let today = new Date();
                this.updated = (today.getMonth() + 1) + "/" + today.getDate() + " " + today.toLocaleTimeString('en-US');
            },
            sortApps: function (apps) { 
                return obj.copy(apps).sort(function(a,b) {
                    if (a.name < b.name)
                        return -1;
                    if ( a.name > b.name)
                        return 1;
                    return 0;
                });
            },
            sortInstances: function (instances) { 
                return obj.copy(instances).sort(function(a,b) {
                    return a.order - b.order;
                });
            },
            refreshApps: function() {
                let that = this

                fetch('/v1/apps')
                .then(function(response) {
                    return response.json()
                })
                .then(function(data) {
                    if (data.message) {
                        return Promise.reject(new Error(data.message))
                    } 
                    
                    that.apps = data;
                    that.updateTime();
                    
                    that.isLoading = false;
                }).catch(function(e) {
                    that.alerts.push({ error: e });

                    that.isLoading = false;
                }); 
            },
            getVersion: function() {
                let that = this

                fetch('/ping')
                .then(function(response) {
                    return response.json()
                })
                .then(function(data) {
                    if (data.message) {
                        return Promise.reject(new Error(data.message))
                    }

                    that.version = data.version
                }).catch(function(e) {
                    that.alerts.push({ error: e });
                });    
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
            getUser: function() {
                return fetch('/v1/user').then(function(response) {
                    return response.json()
                }).then(function(data) {
                    if (data.message) {
                        return Promise.reject(new Error(data.message))
                    }

                    return Promise.resolve(data);
                })
            },
            handleContinue: function () {
                this.modal.confirm.action()
            },
            stopInstance: function () {
                this.scaleInstance(this.selected.app.name, this.selected.instance.name, 0)
            },
            restartInstance: function () {
                this.forceRestartInstance(this.selected.app.name, this.selected.instance.name, this.selected.instance.desiredCount)
            },
            statusClass: function (app, inst) {

                if (inst.error.length > 0 && inst.errorType != "action") {
                    return 'has-background-danger has-text-white'
                }

                if (inst.deployment.isPending) {
                    return 'has-background-warning has-text-black'
                }

                if (inst.isRunning) {
                    return 'has-background-primary has-text-white'
                }

                if (app.type == 'lambda' || app.type == 's3') {
                    return 'has-background-info has-text-white'
                }
                
                if (app.type == 'scheduled-task') {
                    if (inst.cronEnabled) {
                        return 'has-background-info has-text-white'
                    }
                }

                return 'has-background-grey-light has-text-white'
            },
            actions: function(type) {
                switch (type) {
                    case "s3" :
                        return {}
                    case "lambda":
                        return {
                            start: true,
                        }
                    case "scheduled-task": 
                    case "service":
                        return {
                            start: true,
                            scale: true,
                            restart: true,
                            stop: true,
                        }
                }
                
                return {}
            },
            scaleInstance: function(app, instance, count) {
                let that = this
                
                fetch('/v1/apps/' + app + '/instances/' + instance + '/scale/' + count, {
                    method: "PUT",
                    body: JSON.stringify({
                        "version": this.selected.instance.version
                    }),
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
                    
                    let inst = that.getInstance(app, instance); 
                    inst.deployment.isPending = true;
                    that.updateInstance(app, instance, inst);
                }).catch(function(e) {
                    that.alerts.push({ error: e });
                })
            },
            forceRestartInstance: function(app, instance, count) {
                let that = this
                
                fetch('/v1/apps/' + app + '/instances/' + instance + '/restart/' + count, {
                    method: "PUT",
                    body: JSON.stringify({
                        "version": this.selected.instance.version
                    }),
                    headers: {
                        'Accept': 'application/json',
                        'Content-Type': 'application/json' 
                    }
                })
                .then(function(response) {
                    return response.json()
                })
                .then(function(data) {                
                    let inst = that.getInstance(app, instance); 
                    inst.deployment.isPending = true;
                    that.updateInstance(app, instance, inst);
                })
            },
        }
    })
})

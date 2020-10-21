import {obj} from "../../component/copy/object.js";
import {includeTenplates} from "../../component/html/include.js";

includeTenplates().then(() => {
    var app = new Vue({
        el: '#app',

        data: {
            user: { email: "", admin: false },
            config: { consoleLink: "" },

            apps: [],
            projects: [],

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

            filter: {
                terms: "",
                state: "",
            },

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
            'filter.terms': function (val) {
                localStorage.setItem("filter.terms", val);

                this.filterApps()
            },
            'filter.state': function (val) {
                localStorage.setItem("filter.state", val);

                this.filterApps()
            },
        },

        created: function () {
            this.getVersion()

            this.getConfig()

            let view = localStorage.getItem("view")

            if (view) {
                this.setView(view)
            }

            let filterTerms = localStorage.getItem("filter.terms")
            if (filterTerms) {
                this.filter.terms = filterTerms
            }

            let filterstate = localStorage.getItem("filter.state")
            if (filterstate) {
                this.filter.state = filterstate
            }

            let that = this;
            this.getUser().then(function(user) {
                that.user = user

                that.watchForUpdates()
            }).catch(function(e) {
                that.alerts.push({ error: e });
            }); 

            this.filterApps = this.debounce(this.refreshApps, 500);
        },

        methods: {
            toggleProject(e) {
                this.projects = this.projects.map((p) => {
                    if (p.name === e.name) {
                        p.collapsed = !p.collapsed;
                    }
                    return p;
                });
            },
            collapseProjects() {
                this.projects = this.projects.map((p) => {
                    if (p.apps.length > 1) {
                        p.collapsed = true;
                    }
                    return p;
                });
            },
            jumpTo() {
                let hash = window.location.hash.replace("#", "")
            
                var elmnt = document.getElementById(hash);

                if (elmnt) {
                    elmnt.scrollIntoView();
                }
            },
            debounce(func, wait, immediate) {
                var timeout, result;
                return function() {
                  var context = this, args = arguments;
                  var later = function() {
                    timeout = null;
                    if (!immediate) result = func.apply(context, args);
                  };
                  var callNow = immediate && !timeout;
                  clearTimeout(timeout);
                  timeout = setTimeout(later, wait);
                  if (callNow) result = func.apply(context, args);
                  return result;
                };
            },
            formatErrorPreview(error) {
                if (error.length > 16) {
                    return error.substring(0, 16) + "..."
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
            formatDetails(version, revision) {
                return `${version} (revision ${revision})`
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
                
                app.isRefreshing = true;

                let item = that.updateProject(app);

                Vue.set(that.projects, item.i, item.project);

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

                        that.clearErrorsBy("connection");

                        let jumpToHash = (that.projects.length === 0)

                        that.refreshApps(jumpToHash);
                    };

                    source.onerror = function(e) {
                        if (e.readyState == EventSource.CLOSED) {
                            console.log('closed')
                        }

                        that.addError(new Error("real-time updates connection lost"), "connection");
                    };
                    
                    source.addEventListener('message', function(e) {
                        let app = JSON.parse(e.data);

                        let item = that.updateProject(app);

                        Vue.set(that.projects, item.i, item.project);
                        that.updateTime();
                        
                    }, false);
                } else {
                    that.alerts.push({ error: new Error("browser does not support real-time updates")})
                }
            },
            showProject(proj) {
                for (let i in proj.apps) {
                    if (this.showApp(proj.apps[i])) {
                        return true    
                    }
                }

                return false
            },
            showApp(app) {

                for (let i in app.instances) {
                    if (this.showInstance(app.instances[i])) {
                        return true    
                    }
                }

                return false
            },
            showInstance(inst) {
                if (this.filter.state.length == 0) {
                    return true
                }

                switch(this.filter.state) {
                    case "error":
                        if (inst.error.length > 0) {
                            return true
                        }
                        break;
                    case "pending":
                        return inst.deployment.isPending
                    case "running":
                        return inst.isRunning
                    case "stopped":
                            return !inst.deployment.isPending && !inst.isRunning
                }

                return false
            },
            getInstance(appName, instanceName) {
                for (let p in this.projects) {
                    for (let a in this.projects[p].apps) {
                        if (this.projects[p].apps[a].name == appName) {
                            for (let i in this.projects[p].apps[a].instances) {
                                if (this.projects[p].apps[a].instances[i].name == instanceName) {
                                    return this.projects[p].apps[a].instances[i];
                                }
                            }
                        }
                    }
                }

                return undefined
            },
            updateInstance(appName, instanceName, instance) {
                for (let p in this.projects) {
                    for (let a in p.apps) {
                        if (p.apps[a].name == appName) {
                            for (let i in p.apps[a].instances) {
                                if (p.apps[a].instances[i].name == instanceName) {
                                    p.apps[a].instances[i] = instance;
                                }
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
            refreshApps: function(jump) {
                let that = this

                this.isLoading = true;

                let body = {
                    terms: this.filter.terms.length > 0 ? this.filter.terms.trim().split(" ") : [],
                }

                fetch('/v1/apps/filter', {
                    method: "POST",
                    body: JSON.stringify(body),
                    headers: {
                        'Accept': 'application/json',
                        'Content-Type': 'application/json' 
                    }
                }).then(function(response) {
                    return response.json()
                }).then(function(data) {
                    if (data.message) {
                        return Promise.reject(new Error(data.message))
                    } 

                    that.projects = that.groupByProject(data)
                    
                    that.updateTime();
                }).catch(function(e) {

                    that.alerts.push({ error: e });                    
                }).finally(function() {
                    that.isLoading = false;

                    if (jump) {
                        that.jumpTo();
                    }
                }); 
            },
            updateProject: function(app) {
                for (let x=0; x < this.projects.length; x++) {
                    let proj = this.projects[x];

                    if (proj.id === this.extractProject(app).id) {
                        for (let y=0; y < proj.apps.length; y++) {
                            if (proj.apps[y].id === app.id) {
                                proj.apps[y] = app;

                                return { i: x, project: proj };
                            }
                        }
                    }
                }

                return { i: -1, project: {} };
            },
            extractProject: function(app) {
                const noProjectId = "000000000000000000000000";

                return (app.project.id === noProjectId) 
                    ? { id: `${noProjectId}-${app.id}`, name: app.name, is: false }
                    : Object.assign(app.project, { is: true });
            },
            groupByProject: function(apps) {
                let projects = {};

                for (let x=0; x < apps.length; x++) {
                    let proj = this.extractProject(apps[x]);
                    
                    if (proj.id in projects) {
                        let p = projects[proj.id];

                        p.apps.push(apps[x]);
                        
                        projects[proj.id] = p;
                    } else {
                        projects[proj.id] = {
                            id: proj.id,
                            name: proj.name,
                            apps: [apps[x]],
                            is: proj.is
                        }
                    }
                }
               
                let list = [];
                for (let [key, value] of Object.entries(projects)) {
                    list.push(value)
                }

                return list
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
            actions: function(type, inst) {
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
                            scale: !inst.autoScale,
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

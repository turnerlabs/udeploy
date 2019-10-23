import {url} from "../../component/url/params.js";
import {obj} from "../../component/copy/object.js";
import {includeTenplates} from "../../component/html/include.js";

includeTenplates().then(() => {
    var instance = new Vue({
        el: '#instance',

        data: {
            params: url.current.params("/apps/:app/instance/:instance"),
            
            user: {email:""},
            config: {consoleLink: ""},

            alerts: [],
            commits: [],

            app: {},
            instance: {},

            audit: {
                entries: []
            },

            override: {
                env: false
            },

            isLoading: false,
            isPartialLoading: false,
            isLoadingCommits: false,

            page: {},
            updated: ""
        },

        created: function () {
            this.getUser();

            this.getConfig();

            this.refresh(true);

            this.watchForUpdates();

            this.getAudit();

            this.getCommits();
        },

        methods: {
            addError(err) {
                for (let a in this.alerts) {
                    if (this.alerts[a].error.message == err.message) {
                        return
                    } 
                } 

                this.alerts.push({ error: err });
            },
            sortTasks: function (tasks) { 
                return obj.copy(tasks).sort(function(a,b) {
                    if (new Date(a.lastStatusTime) < new Date(b.lastStatusTime))
                        return 1;
                    if (new Date(a.lastStatusTime) > new Date(b.lastStatusTime))
                        return -1;
                    return 0;
                });
            },
            getCommits() {
                let that = this

                that.isLoadingCommits = true; 

                fetch('/v1/apps/'+this.params.app+"/instances/"+this.params.instance+"/commits")
                .then(function(response) {
                    return response.json()
                })
                .then(function(data) {
                    if (data.message) {
                        return Promise.reject(new Error(data.message))
                    }

                    return Promise.resolve(data);
                })
                .then(function(commits) {
                    that.commits = commits;
                })
                .catch(function(e) {
                    that.alerts.push({ error: e });
                })
                .finally(function() {
                    that.isLoadingCommits = false; 
                });
            },
            getAudit() {
                let that = this
                fetch('/v1/apps/'+this.params.app+"/instances/"+this.params.instance+"/audit")
                .then(function(response) {
                    return response.json()
                })
                .then(function(data) {
                    if (data.message) {
                        return Promise.reject(new Error(data.message))
                    }

                    return Promise.resolve(data);
                })
                .then(function(entries) {
                    that.audit.entries = entries;
                })
                .catch(function(e) {
                    that.alerts.push({ error: e });
                });
            },
            refresh(showPageLoading) {
                this.isPartialLoading = true;
                
                if (showPageLoading) {
                    this.isLoading = true 
                }
                
                let that = this;

                fetch('/v1/apps/'+this.params.app)
                .then(function(response) {
                    return response.json()
                })
                .then(function(data) {
                    if (data.message) {
                        return Promise.reject(new Error(data.message))
                    }

                    return Promise.resolve(data);
                })
                .then(function(app) {
                    that.updateInstance(app);

                    that.updateTime();
                    
                }).catch(function(e) {
                    that.addError(e);
                }).finally(function() {
                    that.isPartialLoading = false;
                    that.isLoading = false
                }); 
            },
            updateInstance: function (app) {
                let that = this;

                this.app = app
                
                for (let i in app.instances) {
                    if (app.instances[i].name == this.params.instance) {
                        this.instance = app.instances[i];   
                    }
                }

                if (this.instance.error && this.instance.error.length > 0) {
                    this.addError(new Error(this.instance.error));
                }

                for (let l in this.instance.links) {
                    Object.keys(this.instance.tokens).forEach(function(key) {
                        that.instance.links[l].url = that.instance.links[l].url.replace("{{" + key + "}}", that.instance.tokens[key]);
                    })
                }
            },
            updateTime: function () {
                let today = new Date();
                this.updated = (today.getMonth() + 1) + "/" + today.getDate() + " " + today.toLocaleTimeString('en-US');
            },
            watchForUpdates() {
                let that = this;

                if (!!window.EventSource) {
                    var source = new EventSource('/events/app/changes');

                    source.onopen = function(){
                        console.log('connected');
                    };

                    source.onerror = function(e) {
                        if (e.readyState == EventSource.CLOSED) {
                            console.log('closed')
                        }
                    };
                    
                    source.addEventListener('message', function(e) {
                        let app = JSON.parse(e.data);

                        if (that.app.name == app.name) {
                            that.updateInstance(app);
                        }
                       
                    }, false);
                } else {
                    that.error = "This browser does not support real-time updates. Refresh browser to view changes."
                    // Result to xhr polling :(
                }
            },
            formatStatus: function(status, time) {
                switch (status) {
                    case "STOPPED":
                        return status + " at " + this.formatTime(time) + ""
                    case "RUNNING":
                        return status + " since " + this.formatTime(time) + ""
                    default:
                        return status
                }
            },
            formatCron: function(cron) {
                return cronstrue.toString(cron);
            },
            formatTime: function (today) {
                return (today.getMonth() + 1) + "/" + today.getDate() + "/" + today.getFullYear() + " " + today.toLocaleTimeString('en-US');
            },
            formatEnv(env) {
                let envFile = ""

                for (key in env) {
                    envFile = `${envFile}${key}=${env[key]}\n`
                }

                return envFile
            },
            sortEntries: function (entries) { 
                let temp = obj.copy(entries)
                
                return temp.sort(function(a,b) {
                    return a.time - b.time;
                });
            },
            getUser: function() {
                let that = this

                fetch('/v1/user')
                .then(function(response) {
                    return response.json()
                })
                .then(function(data) {
                    if (data.message) {
                        return Promise.reject(new Error(data.message))
                    }

                    return Promise.resolve(data);
                })
                .then(function(user) {
                    that.user = user
                }).catch(function(e) {
                    that.errors.push({ error: e });
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
        }
    })
});

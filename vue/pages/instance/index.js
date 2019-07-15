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

            this.refresh();

            this.watchForUpdates();

            this.getAudit();

            this.getCommits();
        },

        methods: {
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
            refresh() {
                this.isPartialLoading = true;
                this.isLoading = true 
            
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
                    that.app = app

                    for (let i in app.instances) {
                        if (app.instances[i].name == that.params.instance) {
                            that.instance = app.instances[i]

                            if (that.instance.error && that.instance.error.length > 0) {
                                that.alerts.push({ error: new Error(that.instance.error) })
                            }
                        }
                    }

                    for (let l in that.instance.links) {
                        Object.keys(that.instance.tokens).forEach(function(key) {
                            that.instance.links[l].url = that.instance.links[l].url.replace("{{" + key + "}}", that.instance.tokens[key]);
                        })
                    }

                    that.updateTime();
                    
                }).catch(function(e) {
                    that.alerts.push({ error: e });
                }).finally(function() {
                    that.isPartialLoading = false;
                    that.isLoading = false
                }); 
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
                            that.refresh()
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

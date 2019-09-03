import {env} from "../../env/parse.js";
import {obj} from "../../copy/object.js";

Vue.component('deploy-modal', {
    template: '#deploy-modal-template',
    
    props: ['instance', 'app', 'source'],

    data: function() {
        let data = {
            versions: {
                data: {},
                isLoading: false,
            },
            
            selectedSource: "",
            selectedVersion: "",
            baseVersion: "",

            override: {
                env: false
            },

            environment: "",
            secrets: "",

            deploying: false,
            deployAuthorized: (this.instance.deployCode.length == 0),

            commits: [],
            showChanges: false,
            loadingChanges: false,

            error: "",
            warn: "",
        }

        if (this.instance.containers.length > 0) {
            data.environment = this.formatEnv(this.instance.containers[0].environment)
            data.secrets = this.formatEnv(this.instance.containers[0].secrets)
        }

        return data
    },

    watch: {
        selectedVersion: function(key) {
            let build = this.versions.data[key];

            this.baseVersion = build.version

            this.getCommits(build.version);
        },
        selectedSource: function (instance) {
            let that = this;

            this.versions.isLoading = true;

            this.error = "";
            this.warn = "";
            this.commits = [];

            this.loadVersions(instance)
            .then(function(versions) {
                if (versions.message) {
                    that.error = versions.message;
                } 

                that.versions.data = {}

                Object.keys(versions).map(function(key, index) {
                    let build = versions[key];

                    build.display = that.formatVersion(key, build)
                    
                    that.versions.data[key] = build;
                });

                for (let i in that.app.instances) {
                    if (that.selectedSource == that.app.instances[i].name) {
                        if (that.versions.data[that.app.instances[i].formattedVersion]) {
                            that.selectedVersion = that.app.instances[i].formattedVersion;
                        } else {
                            that.error = that.app.instances[i].formattedVersion + " not found in " + that.selectedSource + " registry"
                        }
                    }
                }

                that.versions.isLoading = false
            })
        }
    },

    mounted: function () { 
        let registry = (this.source)
            ? this.source
            : (this.instance.task.registry && this.instance.task.registry.length > 0) 
                ? this.instance.task.registry
                : this.instance.name

        for (let i in this.app.instances) {
            if (this.app.instances[i].name == registry) {
                this.selectedSource = this.app.instances[i].name
                this.selectedVersion = this.app.instances[i].formattedVersion;

                this.versions.data[this.selectedVersion] = {
                    version: this.app.instances[i].version,
                    display: this.selectedVersion
                }
            }
        }
    },

    methods: {
        sortInstances: function (instances) {
            let temp = obj.copy(instances)
            
            return temp.sort(function(a,b) {
                return a.order - b.order;
            });
        },
        getCommits(version) {
            this.commits = [];
            this.loadingChanges = true

            let that = this
            fetch('/v1/apps/'+this.app.name+"/version/range/"+this.instance.version+"/to/"+version+"/commits")
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
                that.warn = e.message;
            })
            .finally(function() {
                that.loadingChanges = false
            });
        },
        deploy: function (instance) {
            this.error = ""
            this.deploying = true
           
            let that = this

            let body = {
                "version": this.selectedVersion
            }

            if (this.override.env) {
                body.override = true
                body.env = env.parse(this.environment, {});
                body.secrets = env.parse(this.secrets, {});
            }

            let ver = this.versions.data[this.selectedVersion]

            if (ver.registry) {
                body.imageTag = this.selectedVersion;
            }

            return fetch('/v1/apps/' + this.app.name + "/instances/" + instance + "/deploy/" + this.selectedSource + "/" + ver.revision, {
                method: "POST",
                body: JSON.stringify(body),
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
                    that.error = data.message;
                    that.deploying = false
                } else {
                    that.$parent.modal.deploy.show = false
                }
            })
        },
        loadVersions: function (instance) {
            return fetch('/v1/apps/' + this.app.name + "/instances/" + instance + "/registry")
            .then(function(response) {
                return response.json()
            })
        },
        validateDeployCode: function (evt) {
            this.deployAuthorized = (this.instance.deployCode == evt.target.value)
        },
        overrideEnv: function(evt) {
            this.override.env = evt.target.checked;
        },
        formatVersion(version, build) {
            if (build.registry) {
                return version;
            }

            return version + " (" + build.revision + ")"; 
        },
        formatEnv(env) {
            let envFile = ""

            for (let key in env) {
                envFile = `${envFile}${key}=${env[key]}\n`
            }

            return envFile
        },
        willPropagate() {
            for(let i in this.app.instances) {
                if (this.app.instances[i].task.registry == this.instance.name && this.app.instances[i].autoPropagate) {
                    return true
                }
            }

            return false
        }
    }
})
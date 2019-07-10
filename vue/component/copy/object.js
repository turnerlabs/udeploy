export let obj = {
    copy: function (o) {
        var output, v, key;
        output = Array.isArray(o) ? [] : {};
        for (let key in o) {
            v = o[key];
            output[key] = (typeof v === "object") ? this.copy(v) : v;
        }
        return output;
    },
}


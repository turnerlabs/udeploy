export let formatCommitMessage = (msg) => {
    msg = msg.replace(/\n/gi, "<br/>")
    return msg
}

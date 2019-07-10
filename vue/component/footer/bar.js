Vue.component('footer-bar', {
    props: ['updated', 'version'],

    template: `<div class="field is-grouped is-grouped-multiline is-pulled-right">
    <div v-if="updated.length > 0" class="control">
        <span class="tags has-addons">
            <span class="tag is-dark">updated</span>
            <span class="tag is-primary">{{ updated }}</span>
        </span>
    </div>
    <div v-if="version" class="control">
        <span class="tags has-addons">
            <span class="tag is-dark">version</span>
            <span class="tag is-primary">{{ version }}</span>
        </span>
    </div>
</div>`
})
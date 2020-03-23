Vue.component('header-bar', {
    props: ['user', 'page', 'config'],

    template: `
<div class="nav">
<div class="container">
    <div class="is-pulled-right options">
        <table>
            <tr>
                <td>
                    <span class="buttons has-addons">
                        <a class="button is-primary is-small" v-bind:class="{ 'is-outlined': !page.apps }" :disabled="page.apps" href="/apps">Apps</a>
                        <a v-if="user.admin" class="button is-primary is-small" v-bind:class="{ 'is-outlined': !page.projects }" :disabled="page.projects" href="/apps/projects">Projects</a>
                        <a v-if="user.admin" class="button is-primary is-small" v-bind:class="{ 'is-outlined': !page.notices }" :disabled="page.notices" href="/apps/notices">Notices</a>
                        <a v-if="user.admin" class="button is-primary is-small" v-bind:class="{ 'is-outlined': !page.users }" :disabled="page.users" href="/apps/users">Users</a>
                    </span>
                </td>
                <td class="signout">
                    <a class="button is-dark is-outlined is-small" v-bind:href="config.consoleLink" target="_console">AWS Console</a>
                    <a class="button is-dark is-outlined is-small" href="/oauth2/logout" title="Clearing browser site data may be required when signing out.">Sign Out</a>
                </td>
            </tr>
        </table>
        <div class="email is-pulled-right has-text-grey-light">{{ user.email }}</div>
    </div>
    <p class="has-text-weight-bold is-size-2"><a class="has-text-black" href="/apps">uDeploy</a></p>
    A simple way to deploy AWS resources.
</div>
</div>`
})
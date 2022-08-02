<template>
  <q-page class="flex flex-center">

    <div v-if="this.authenticated === true">
      <h3 class="justify-center">Callback Page</h3>
      <div>You are now connected, you can go to the Home page</div>
    </div>
    <div v-else>
      <h3 class="justify-center">Callback Page</h3>
      <div>You are not connected, go to the Login page</div>
    </div>

  </q-page>
</template>

<script>
import { defineComponent } from 'vue'
import { getAuthenticated, loadLoginScreen, login, logout, validateState } from '../api/oidc';

export default defineComponent({
  name: 'CallbackPage',
  methods: {

  },
  data() {
    return {
      authenticated: false
    }
  },
  async created() {
    console.log('CallbackPage just created');
    const params = (new URL(document.location)).searchParams;
    const state = params.get('state');
    const code = params.get('code');
    const auth = await getAuthenticated();
    console.log("authenticated: ", auth);
    if (auth) {
      this.authenticated = true;
    } else {
      if (code !== null) {
        try {
            const validState = validateState(state);
            if (!validState) throw new Error();
            this.authenticated = await login(code);
        } catch (err) {
            // DO NOTHING
        }
      }
    }
  },
  })

</script>

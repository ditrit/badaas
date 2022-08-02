<template>
  <q-page class="flex flex-center">
    
    

    <div v-if="this.authenticated === true">
      <h3 class="justify-center">Home page</h3>
      <div>You are connected</div>
    </div>

    <div v-else>
      <h3 class="justify-center">Home page</h3>
      <div>You are not connected</div>
    </div>

  </q-page>
</template>

<script>
import { defineComponent } from 'vue'
import { getAuthenticated, loadLoginScreen, login, logout, validateState } from '../api/oidc';

export default defineComponent({
  name: 'HomePage',
  methods: {

  },
  data() {
    return {
      authenticated: false
    }
  },
  async created() {
    console.log('HomePage just created');
    const auth = await getAuthenticated();
    console.log("authenticated:", auth);
    if (auth) {
      this.authenticated = true
    }
  },
  })

</script>

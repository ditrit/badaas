<template>
  <q-page class="flex flex-center">
    <!-- <img
      alt="Quasar logo"
      src="~assets/quasar-logo-vertical.svg"
      style="width: 200px; height: 200px"
    > -->
    <div v-if="this.authenticated === false">
      <h3 class="justify-center">Login/Logout Page</h3>
      <q-btn @click="() => LoginClick('google')" color="primary" icon="login" label="Login with Google" />
      &nbsp;
      <q-btn @click="() => LoginClick('gitlab')" color="primary" icon="login" label="Login with Gitlab" />
    </div>
    <div v-else>
      <h3 class="justify-center">Login/Logout Page</h3>
      <div>Hello {{this.user}} !</div>
      <q-btn @click="LogoutClick" color="primary" icon="logout" label="Logout" />
    </div>

    
  </q-page>
</template>

<script>
import { defineComponent } from 'vue'
import { getAuthenticated, loadLoginScreen, login, logout, validateState, revokeToken, logoutSession } from '../api/oidc';
import jwt_decode from 'jwt-decode';

export default defineComponent({
  name: 'HomePage',
  methods: {
    LoginClick(prov){
      window.localStorage.setItem('provider', prov);
      loadLoginScreen();
    },
    async LogoutClick(){
      this.authenticated = false;
      logoutSession();
    }
  },
  data() {
    return {
      authenticated: false,
      user: ""
    }
  },
  async created() {
    console.log('IndexPage just created');
    const auth = await getAuthenticated();
    console.log("authenticated:", auth);
    if (auth) {
      this.authenticated = true;
      const email = "RANDOM USER";
      this.user = email;
    }
  },
  })

</script>

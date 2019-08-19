<template>
  <section>
    <div id="firebaseui-auth"></div>
  </section>
</template>

<script> 
import * as firebase from 'firebase/app'
import * as firebaseui from 'firebaseui'

import 'firebaseui/dist/firebaseui.css'

var ui = null

export default {
  name: 'Auth',
  props: {
    msg: String
  },
  created () {
    ui = new firebaseui.auth.AuthUI(firebase.auth())
  },
  mounted () {
    ui.start('#firebaseui-auth', {
      signInSuccessUrl: '/',
      signInOptions: [
        firebase.auth.GithubAuthProvider.PROVIDER_ID,
        firebase.auth.TwitterAuthProvider.PROVIDER_ID,
        firebase.auth.GoogleAuthProvider.PROVIDER_ID,
        firebase.auth.EmailAuthProvider.PROVIDER_ID
      ],
      credentialHelper: firebaseui.auth.CredentialHelper.GOOGLE_YOLO
    });
  }
}
</script>

<style scoped>
</style>

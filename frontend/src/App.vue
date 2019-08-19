<template>
  <div id="app">
    <b-navbar type="is-light">
        <template slot="brand">
            <b-navbar-item type="div">
                <img alt="Bostadskoe logo" src="./assets/logo.png">
            </b-navbar-item>
        </template>

        <template slot="end">
            <b-navbar-item href="https://github.com/phille97/bostadskoe" target="_BLANK">
              <b-icon
                icon="github-circle"
                size="is-medium" />
            </b-navbar-item>
            <b-navbar-item tag="div" v-if="user">
                <div class="buttons">
                    <b-button @click="signOut">Avbryt</b-button>
                </div>
            </b-navbar-item>
        </template>
    </b-navbar>
    

    <section class="section">
    <div class="container">
    <b-steps v-model="activeStep" :animated="true" :has-navigation="false" size="is-large">
        <b-step-item icon="account-key" :clickable="true" :type="{'is-success': !!user}">
          <Auth />
        </b-step-item>

        <b-step-item icon="filter" :clickable="!!user">
          <Form />
          <h1>TODO</h1>
        </b-step-item>

        <b-step-item icon="bell" :clickable="!!user">
          <h1>TODO</h1>
        </b-step-item>

        <b-step-item icon="check" :clickable="false">
          <h1>🎉🎉🎉 Nu kommer du få e-mail om något nytt dyker upp 😃</h1>
        </b-step-item>
    </b-steps>
      </div>
    </section>
  </div>
</template>

<script>
import * as firebase from 'firebase/app'

import Form from './components/Form.vue'
import Auth from './components/Auth.vue'

import 'buefy/dist/buefy.css'

export default {
  name: 'app',
  data () {
    return {
      activeStep: 0,
      user: firebase.auth().currentUser
    }
  },
  components: {
    Form,
    Auth
  },
  methods: {
    signOut () {
      firebase.auth().signOut();
      this.activeStep = 0
    }
  },
  created () {
      this.user = firebase.auth().currentUser
  },
  mounted () {
    firebase.auth().onAuthStateChanged(
      user => {
        this.user = user
        this.activeStep = this.user ? 1 : 0
      },
      error => {
        console.log(error);
      }
    )
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>

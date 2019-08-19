import Vue from 'vue'
import Buefy from 'buefy'
import App from './App.vue'
import * as firebase from 'firebase/app'
import 'firebase/auth'

Vue.config.productionTip = false

const firebaseConfig = {
    apiKey: "****",
    authDomain: "bostadskoe.firebaseapp.com",
    databaseURL: "https://bostadskoe.firebaseio.com",
    projectId: "bostadskoe",
    storageBucket: "bostadskoe.appspot.com",
    messagingSenderId: "901462286394",
    appId: "1:901462286394:web:5cf7ae5550f2c03c"
}
firebase.initializeApp(firebaseConfig)

Vue.use(Buefy)

new Vue({
  render: h => h(App),
}).$mount('#app')

<template>
  <div class="hello">
    <h1>{{ title }}</h1>
    <ul>
        <!--<li v-for="residence in residences">
            <pre>{{ residence }}</pre>
        </li>-->
    </ul>
  </div>
</template>

<script>
import { db } from '../db'
export default {
  name: 'ResidenceList',
  props: {
    title: String,
    collection: String
  },
  data () {
    return {
      residences: []
    }
  },
  watch: {
    collection: {
      immediate: true,
      handler (collection) {
        this.$bind('residences', db.collection(collection))
      }
    }
  }
}
</script>

<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
li.synced a {
  color: #42b983;
}
li a {
  color: #ff8579;
}
</style>

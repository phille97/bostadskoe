<template>
  <div class="hello">
    <h1>Bostadskö</h1>
    <!--<hr>
        <ResidenceList title="Senäte" collection="senate" />
        <ResidenceList title="Bostadsförmedlingen i Stockholm" collection="bostadstockholm" />
        <hr>-->
  </div>
</template>

<script>
import ResidenceList from './ResidenceList'

export default {
  name: 'Index',
  components: { ResidenceList }
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

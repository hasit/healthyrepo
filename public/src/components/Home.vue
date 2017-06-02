<template>
  <div class="healthyrepo">
    <h1>{{ title }}</h1>
    <h2>See Health Indicators of Github Repositories</h2>

    <span>github.com/<input type="text" placeholder="owner" v-autosize="repoOwner">{{repoOwner}}</input>/<input type="text" placeholder="repo" v-autosize="repoName">{{repoName}}</input> <a v-on:click="">⏎</a> </span>
  </div>
</template>

<script>
export default {
  name: 'healthyrepo',
  data () {
    return {
      title: 'HealthyRepo'
    }
  },
  created () {
    this.$http
    .get('https://api.healthyrepo.com/indicators', (data) => {
      console.log(data)
    })
    .error((err) => console.log(err))
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1 {
  font-weight: 900;
}

h2 {
  font-weight: 600;
}

span {
  font-size: 2em;
  font-weight: 300;
}

input {
  border: 0;
  outline: 0;
  background: transparent;
  font-size: 1em;
  font-weight: 300;
  color: #294455;
  text-align: center;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #D0031E;
  text-decoration: none;
   cursor: pointer;
}
</style>

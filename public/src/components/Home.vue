<template>
  <div class="healthyrepo">
    <h1><a href="/">{{ title }}</a></h1>

    <ul v-if="errors && errors.length">
      <li v-for="error of errors" class="error">
        {{error.message}}
      </li>
    </ul>

    <h2>1. Find Github repository</h2>
    <span class="inputs" @keyup.enter="getRepo">github.com/<input type="text" placeholder="owner" v-autosize="repoOwner" v-model="owner">{{repoOwner}}</input>/<input type="text" placeholder="repo" v-autosize="repoName" v-model="repo">{{repoName}}</input></span>

    <div v-if="showRepo">
      <a v-bind:href="repository.owner_url" target="_blank"><img v-bind:src="repository.owner_avatar_url" height="64" width="64"/></a>
      <br>
      <a v-bind:href="repository.repo_url" target="_blank">{{ repository.full_name }}</a>

      <h2>2. Click indicator of interest</h2>
      <ul class="indicators">
        <li class="indicator" v-for="indicator in indicators" v-on:click="getIndicator(indicator.key)" v-if="indicator.active" v-bind:title="indicator.description">
          {{ indicator.name }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'healthyrepo',
  data () {
    return {
      owner: '',
      repo: '',
      repoOwner: '',
      repoName: '',
      title: 'HealthyRepo',
      repository: {},
      showRepo: false,
      indicators: [],
      errors: []
    }
  },
  methods: {
    getRepo: function () {
      this.errors = []
      this.showRepo = false

      if (this.owner === '' || this.repo === '') {
        var e = new Error('Fields \'owner\' or \'repo\' cannot be empty')
        this.errors.push(e)
        return
      }

      var reqURL = 'http://localhost:8080/repo/' + this.owner + '/' + this.repo
      console.log(reqURL)
      axios({
        method: 'get',
        url: reqURL,
        headers: {
          'Access-Control-Allow-Origin': '*'
        }
      })
      .then(response => {
        console.log(response.data)
        this.repository = response.data
        this.showRepo = true
      })
      .then(() => {
        var reqURL = 'http://localhost:8080/indicators'
        console.log(reqURL)
        axios({
          method: 'get',
          url: reqURL,
          headers: {
            'Access-Control-Allow-Origin': '*'
          }
        })
        .then(response => {
          console.log(response.data)
          this.indicators = response.data
        })
        .catch(e => {
          this.errors.push(e)
        })
      })
      .catch(e => {
        console.log(e)
        this.errors.push(e)
      })
    },
    getIndicator: function (key) {
      var reqURL = 'http://localhost:8080/repo/' + this.owner + '/' + this.repo + '/health/' + key
      console.log(reqURL)
      axios({
        method: 'get',
        url: reqURL,
        headers: {
          'Access-Control-Allow-Origin': '*'
        }
      })
      .then(response => {
        console.log(response.data)
      })
      .catch(e => {
        this.errors.push(e)
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1 {
  font-weight: 900;
  margin-top: 0;
}

h2 {
  font-weight: 600;
  margin: 0;
}

span.inputs{
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
  margin-top: 0;
}

li.error {
  display: inline-block;
  margin: 0 10px;
  cursor: default;
  color: #D0031E;
}

li.indicator {
  display: inline-block;
  margin: 0 10px;
  cursor: default;
  font-size: 2em;
  border-bottom: 1px dotted #294455;
  font-weight: 300;
}

li.indicator:hover {
  display: inline-block;
  margin: 0 10px;
  cursor: pointer;
  color: #D0031E;
}

a {
  color: #294455;
  text-decoration: none;
  cursor: pointer;
}

a:hover {
  color: #D0031E;
  text-decoration: none;
  cursor: pointer;
}

img {
  margin-top: 5px;
}

</style>

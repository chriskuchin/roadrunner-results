<template>
  <table class="table">
    <thead>
      <tr>
        <th v-for="header in headers">
          <abbr :title="header.title">
            {{ header.abbr !== '' ? header.abbr : header.title }}
          </abbr>
        </th>
      </tr>
    </thead>
    <tfoot>
      <tr>
        <th v-for="header in headers">
          <abbr :title="header.title">
            {{ header.abbr !== '' ? header.abbr : header.title }}
          </abbr>
        </th>
      </tr>
    </tfoot>
    <tbody>
      <tr v-for="row in rows">
        <th v-for="col in row">{{ col }}</th>
      </tr>
    </tbody>
  </table>
</template>

<script>
export default {
  props: ['headers', 'rows'],
  data: function () {
    return {}
  }
}
</script>
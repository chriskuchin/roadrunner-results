<template>
  <div class="container">
    <h1>Lane Assignment</h1>
    <div class="tabs">
      <ul>
        <li class="is-active"><a>Heat 1</a></li>
        <li><a>Heat 2</a></li>
        <li><a>Heat 3</a></li>
        <li><a>Heat 4</a></li>
        <li><a>New +</a></li>
      </ul>
    </div>
    <div class="section">
      <table class="table">
        <thead>
          <tr>
            <th><abbr title="Lane Assignment">Lane</abbr></th>
            <th><abbr title="Athlete Bib Number">Bib</abbr></th>
            <th><abbr title="Athlete First Name">F. Name</abbr></th>
            <th><abbr title="Athlete Last Name">L. Name</abbr></th>
            <th><abbr title="Athlete Birth Year">B. Year</abbr></th>
          </tr>
        </thead>
        <tfoot>
          <tr>
            <th><abbr title="Lane Assignment">Lane</abbr></th>
            <th><abbr title="Athlete Bib Number">Bib</abbr></th>
            <th><abbr title="Athlete First Name">F. Name</abbr></th>
            <th><abbr title="Athlete Last Name">L. Name</abbr></th>
            <th><abbr title="Athlete Birth Year">B. Year</abbr></th>
          </tr>
        </tfoot>
        <tbody>
          <tr v-for="(assignment, index) in lanes">
            <th>{{ assignment.lane }}</th>
            <td><input :tabindex="index + 1" class="input is-small" type="text" placeholder="Text input"
                :value="assignment.bib"></td>
            <td></td>
            <td></td>
            <td></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>

export default {
  mounted: function () {
  },
  data: function () {
    return {
      id: "123",
      lanes: [
        {
          lane: 1,
          bib: "211"
        },
        {
          lane: 2,
          bib: "221"
        },
        {
          lane: 3,
          bib: "272"
        },
        {
          lane: 4,
          bib: "222"
        },
        {
          lane: 5,
          bib: "233"
        },
        {
          lane: 6,
          bib: "243"
        },
        {
          lane: 7,
          bib: "289"
        },
        {
          lane: 8,
          bib: "227"
        },
      ],
    }
  },
  computed: {
  },
  methods: {
  }
};
</script>
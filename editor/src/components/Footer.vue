<template>
  <footer>
    <p>
      &copy; 2010&ndash;{{ year }}
      <a href="https://www.facebook.com/kkeutsori/">Brian</a>
      &amp;
      <a href="https://subl.ee/">Heungsub</a>.
      All rights reserved.
    </p>
  </footer>
</template>

<script>
export default {
  name: 'Footer',

  data: () => ({
    year: (new Date()).getFullYear()
  })
}
</script>

<style scoped>
p {
  text-align: center;
}

.logo {
  padding: 3rem;
  padding-bottom: 0.5rem;
}

.intro {
  padding-bottom: 0.5rem;
}
</style>

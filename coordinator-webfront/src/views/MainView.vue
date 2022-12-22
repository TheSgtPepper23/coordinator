<script setup>
import Sidebar from '../components/Sidebar.vue';
import MapContent from '../components/MapContent.vue';

</script>

<template>
  <v-row>
    <v-col cols="3">
      <Sidebar></Sidebar>
    </v-col>
    <v-col cols="9">
      <MapContent></MapContent>
    </v-col>
  </v-row>
</template>

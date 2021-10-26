<template>
  <!-- 在你的template代码中使用地图系列的组件 -->
  <div class="map-container">
    <amap :zoom="11" :center="[117.000923, 36.675807]">
      <amap-marker :position="[117.000923, 36.675807]" />
    </amap>
  </div>
</template>
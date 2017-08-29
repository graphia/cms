<template>
    <svg :height="height" :class="cssClasses" :viewBox="viewBox" version="1.1" aria-hidden="true" role="presentation" :width="width" :icon-name="iconName" v-html="iconPath">
    </svg>
</template>

<script lang="babel">
  import Octicons from 'octicons';
  import Vue from 'vue';

  export default {
    name: 'octicon',
    props: {
      iconName: {
        type: String,
        validator(value) {
          if (value) {
            if (Octicons[value]) {
              return true;
            }
            Vue.util.warn(`${value} is not a valid octicon icon-name.`, this);
            return false;
          }
          Vue.util.warn('You must specify an octicon icon-name.', this);
          return null;
        },
      },
      viewBox: {
        type: String,
        default() {
          if (Octicons[this.iconName]) {
            return `0 0 ${Octicons[this.iconName].width} ${Octicons[this.iconName].height}`;
          }
          return '0 0 12 16';
        },
      },
    },
    data() {
      return {
        iconHeight: 0,
        iconWidth: 0,
      };
    },
    mounted() {
      if (this.$el) {
        this.iconHeight = parseFloat(window.getComputedStyle(this.$el, null).getPropertyValue('font-size'));
        if (Octicons[this.iconName]) {
          this.iconWidth = (Octicons[this.iconName].width / Octicons[this.iconName].height)
            * this.iconHeight;
        }
      }
    },
    computed: {
      // Add the default CSS class, the icon CSS class and any extra CSS classes together.
      cssClasses() {
        return `octicon octicon-${this.iconName}`;
      },
      // Get the icon's path tag.
      iconPath() {
        if (Octicons[this.iconName]) {
          return Octicons[this.iconName].path;
        }
        return null;
      },
      // Get the icon height from the font-size of the element.
      height() {
        return this.iconHeight || 0;
      },
      // Get the icon width from the icon size and the height of the element.
      width() {
        return this.iconWidth || 0;
      },
    },
  };
</script>

<style>
    .octicon {
        display: inline-block;
        vertical-align: middle;
        fill: currentColor;
    }

    .h1 .octicon, .h2 .octicon, .h3 .octicon, .h4 .octicon, .h5 .octicon, .h6 .octicon,
    h1 .octicon, h2 .octicon, h3 .octicon, h4 .octicon, h5 .octicon, h6 .octicon {
        vertical-align: bottom;
    }

</style>

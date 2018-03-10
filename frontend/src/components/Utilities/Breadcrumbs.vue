<template>
	<nav>
		<ol class="breadcrumb">


			<!-- if we're currently on this page, don't show a link and add active class -->
			<li v-if="currentPage(breadcrumb)" v-for="(breadcrumb, i) in breadcrumbs" :key="i" class="breadcrumb-item active">{{breadcrumb.text}}</li>
			<li v-else class="breadcrumb-item">
				<router-link :to="{name: breadcrumb.target, params: breadcrumb.params}">{{breadcrumb.text}}</router-link>
			</li>
		</ol>
	</nav>
</template>

<script lang="babel">
	import {HomeBreadcrumb} from "../../javascripts/models/breadcrumb.js";
	export default {
		name: "Breadcrumbs",
		props: [
			"levels"
		],
		computed: {
			breadcrumbs() {
				// Home is always the 'first' breadcrumb, so we needn't
				// specify it elsewhere
				return [HomeBreadcrumb].concat(this.levels);
			}
		},
		methods: {
			currentPage(breadcrumb) {
				if (!breadcrumb) {
					console.error("breadcrumb has no target, ensure :levels are supplied");
					return;
				};
				return this.$router.history.current.name === breadcrumb.target;
			}
		}
	};
</script>
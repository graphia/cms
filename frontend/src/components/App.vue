<template>
	<div id="application">

		<nav id="primary">
			<ul>
				<li class="active">Files</li>
				<li>Activity</li>
			</ul>
		</nav>

		<nav id="secondary">
			<ul>
				<li v-for="document in documents">
					{{ document.filename }}
				</li>
			</ul>
		</nav>

		<div id="router-view-container">
			<div id="router-view">
				<router-view/>
			</div>
		</div>
	</div>
</template>

<script lang="babel">

	export default {
		computed: {
			documents() {
				return this.$store.state.documents;
			}
		}
	};
</script>
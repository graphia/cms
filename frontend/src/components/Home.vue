<template>
	<div>
		<h2>Hello World</h2>

		<button class="btn btn-primary" :class="{'disabled': publishing}" @click="publish">
			{{ publishing ? "Publishing" : "Publish" }}
		</button>

	</div>
</template>

<script lang="babel">
	import Broadcast from '../components/Broadcast';
	import CMSPublisher from '../javascripts/publish.js';

	export default {
		name: "Home",

		data() {
			return {
				publishing: false
			};
		},

		methods: {

			async publish(event) {

				this.publishing = true;

				try {
					console.log("Starting publishing");
					await CMSPublisher.publish();
					console.log(this.publishing);

					// TODO it would be nice to include a hyperlink so the
					// user could immediately see their work in its published
					// state, this would require some reworking of the
					// Broadcast compontent and broadcast.js
					this.$store.state.broadcast.addMessage(
						"success",
						"Success!",
						"Documentation published",
						3
					);
				}
				finally {
					this.publishing = false;
				}
			}
		}
	}
</script>

<style lang="scss" scoped>
</style>

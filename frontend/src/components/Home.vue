<template>
	<div>

		<div class="row">
			<div class="col-md-6">
				<div class="card">
					<div class="card-block">
						<h4 class="card-title">Recent Updates</h4>
					</div>
				</div>
			</div>

			<div class="col-md-6">
				<div class="card">
					<div class="card-block">
						<h4 class="card-title">Statistics</h4>
					</div>
				</div>
			</div>
		</div>

		<div class="row">
			<button class="btn btn-primary" :class="{'disabled': publishing}" @click="publish">
				{{ publishing ? "Publishing" : "Publish" }}
			</button>
		</div>

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

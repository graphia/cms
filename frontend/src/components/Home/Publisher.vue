<template>

	<div class="publisher card bg-light">
		<div class="card-body">

			<div class="alert alert-success" v-if="published">
				Publishing was successful, click <a href="/" target="_blank">here</a>
				to view the finished site.
			</div>
			<div class="alert alert-info" v-else>
				When you are happy with the content, you can release it
				to your audience by publishing it.
			</div>

			<button class="btn btn-lg btn-success" :class="{'disabled': publishing}" @click="publish">
				<octicon :icon-name="'cloud-upload'"></octicon>
				{{ publishing ? "Publishing" : "Publish" }}
			</button>

		</div>
	</div>

</template>

<script lang="babel">
	import CMSPublisher from '../../javascripts/publish.js';
	import checkResponse from '../../javascripts/response.js';

	export default {
		name: "Publisher",
		data() {
			return {
				publishing: false,
				published: false
			};
		},
		methods: {
			async publish(event) {

				this.publishing = true;
				console.debug("Starting publishing");

				try {

					let response = await CMSPublisher.publish();
					let json = await response.json();

					if (!checkResponse(response.status)) {
						console.error("Failed to publish");
						console.error(json.meta);
						throw "publishing failed";
						return;
					};

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

					// once we've published successfully at least once, display
					// a success message
					this.published = true;

				}
				catch(error) {
					console.error(error);
				}
				finally {
					this.publishing = false;
				}
			},
		}
	}
</script>

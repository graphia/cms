<template>
	<button @click="deleteDirectory" class="btn btn-danger btn-sm">
		Delete directory
	</button>
</template>


<script lang="babel">

	import config from '../../../javascripts/config.js';
	import CMSCommit from '../../../javascripts/models/commit.js';
	import checkResponse from '../../../javascripts/response.js';

	import Accessors from '../../Mixins/accessors';

	export default {
		name: "DirectoryDelete",
		mixins: [Accessors],
		methods: {

			async deleteDirectory(event) {
				event.preventDefault();

				let title = this.activeDirectory.title

				let commit = new CMSCommit(
					`Deleting directory ${title}`,
					[], // no files
					[this.activeDirectory]
				);

				let response = await this.activeDirectory.destroy(commit);

				if (!checkResponse(response.status)) {
					console.error("could not delete directory", response);
					return;
				};

				// success

				let json = await response.json();

				await this.$store.commit("setLatestRevision", json.oid);

				this.$store.state.broadcast.addMessage(
					"success",
					"Directory deleted",
					`${title} and its contents have been deleted`,
					3
				);

				this.redirectToHome();

			},

			redirectToHome() {
				this.$router.push({name: 'home'});
			}
		}
	};
</script>

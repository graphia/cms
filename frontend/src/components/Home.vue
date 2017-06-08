<template>
	<div>

		<div class="row">
			<div class="col-md-6">
				<div class="card recent-updates">
					<div class="card-block">
						<h4 class="card-title">Recent Updates</h4>
					</div>

					<ol class="list-group list-group-flush">
						<li class="recent-commit-info list-group-item" v-for="commit in commits">
							<a :href="`/cms/commit/${commit.id}`">
								{{ commit.message || "Empty commit message" }}
							</a>
							<p class="card-text">
								<small>{{ commit.author.Name }} committed 2 minutes ago</small>
							</p>
						</li>
					</ol>
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
	import config from '../javascripts/config.js';

	export default {
		name: "Home",

		created() {
			this.commits = this.getLatestCommits();
		},

		data() {
			return {
				publishing: false,
				commits: []
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
			},

			async getLatestCommits() {
				let path = `${config.api}/recent_commits`;

				try {
					let response = await fetch(path, {mode: "cors", headers: this.$store.state.auth.authHeader()});

					if (response.status != 200) {
						throw("Could not retrieve recent commit sumamry");
					}

					let json = await response.json()

					this.commits = json;


				} catch(err) {
					throw(err);
				}

			}
		}
	}
</script>

<style lang="scss" scoped>
</style>

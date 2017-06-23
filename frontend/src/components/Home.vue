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
							<router-link :to="{name: 'commit', params: {hash: commit.id}}">
								{{ commit.message || "Empty commit message" }}
							</router-link>
							<p class="card-text">
								<small>{{ commit.author.Name }} committed 2 minutes ago</small>
							</p>
						</li>
					</ol>
				</div>
			</div>

			<div class="col-md-6">
				<div class="card statistics">
					<div class="card-block">
						<h4 class="card-title">Statistics</h4>

						<dl class="row">
							<dt class="col-sm-2">Users</dt>
							<dd class="col-sm-10">4</dd>

							<dt class="col-sm-2">Commits</dt>
							<dd class="col-sm-10">38</dd>

							<dt class="col-sm-2">Files</dt>
							<dd class="col-sm-10">12</dd>
						</dl>

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
	import checkResponse from '../javascripts/response.js';

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

					if (!checkResponse(response.status)) {
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

	.card.statistics dl {
		margin-bottom: 0rem;
	}
</style>

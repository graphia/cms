<template>
	<div>

		<Breadcrumbs :levels="breadcrumbs" />

		<h1>Dashboard</h1>

		<DirectorySummary class="directories"/>

		<div class="row mt-4">

			<div class="col col-md-4">

				<div class="card bg-light">
					<div class="card-body">

						<p>
							When you are happy with the content, you can release it
							to your audience by publishing it.
						</p>

						<button class="btn btn-lg btn-success" :class="{'disabled': publishing}" @click="publish">
							<octicon :icon-name="'cloud-upload'"></octicon>
							{{ publishing ? "Publishing" : "Publish" }}
						</button>

					</div>
				</div>

			</div>

			<div class="col-md-5">
				<div class="card recent-updates">
					<div class="card-body">
						<h4 class="card-title">Recent Updates</h4>
					</div>

					<ol class="list-group list-group-flush">
						<li class="recent-commit-info list-group-item" v-for="(commit, i) in commits" :key="i">
							<router-link :to="{name: 'commit', params: {hash: commit.id}}">
								{{ commit.message || "Empty commit message" }}
							</router-link>
							<p class="card-text">
								<small>{{ commit.author.Name }} committed {{ commit.timestamp | time_ago }}</small>
							</p>
						</li>
					</ol>
				</div>
			</div>

			<div class="col-md-3">
				<div class="card statistics">
					<div class="card-body">
						<h4 class="card-title">Statistics</h4>

						<dl class="row">
							<dt class="col-sm-9">Users</dt>
							<dd class="col-sm-3">4</dd>

							<dt class="col-sm-9">Commits</dt>
							<dd class="col-sm-3">38</dd>

							<dt class="col-sm-9">Files</dt>
							<dd class="col-sm-3">12</dd>
						</dl>

					</div>

				</div>
			</div>
		</div>

	</div>
</template>

<script lang="babel">
	import Broadcast from '../components/Broadcast';
	import DirectorySummary from '../components/Directory/Summary';
	import Breadcrumbs from '../components/Utilities/Breadcrumbs';
	import CMSBreadcrumb from '../javascripts/models/breadcrumb.js';

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

		components: {
			DirectorySummary,
			Breadcrumbs
		},

		computed: {
			breadcrumbs() {
				return [];
			}
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

					this.commits = await response.json()

				} catch(err) {
					//throw(err);
					console.error("Couldn't retrieve latest commits")
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

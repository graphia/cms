<template>
	<div v-title="title">

		<Breadcrumbs :levels="breadcrumbs" />

		<h1>Dashboard</h1>

		<div class="row mt-4">

			<div class="col col-md-4">
				<Publisher/>
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

		<DirectorySummary class="directories"/>

	</div>
</template>

<script lang="babel">
	import Broadcast from './Broadcast';
	import DirectorySummary from './Directory/Summary';
	import Breadcrumbs from './Utilities/Breadcrumbs';
	import Publisher from './Home/Publisher';

	import CMSBreadcrumb from '../javascripts/models/breadcrumb.js';
	import config from '../javascripts/config.js';
	import checkResponse from '../javascripts/response.js';

	export default {
		name: "Home",

		created() {
			this.commits = this.getLatestCommits();
		},

		data() {
			return {
				commits: [],
				title: "Graphia CMS"
			};
		},

		components: {
			DirectorySummary,
			Breadcrumbs,
			Publisher
		},

		computed: {
			breadcrumbs() {
				return [];
			}
		},

		methods: {

			async getLatestCommits() {
				let path = `${config.api}/recent_commits`;

				try {
					let response = await fetch(path, {headers: this.$store.state.auth.authHeader()});

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

Given /^that the server is running$/ do
  @server = background_process('../../graphia-cms', load: true).with do |process|


    process.argument '-config', '../../config/cucumber.yml'
    process.argument '-log-to-file', 'true'

    process.ready_when_url_response_status 'http://localhost:9095/health.test', 'OK'

  end
end
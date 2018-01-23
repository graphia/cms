Given %r{^my user account with public key exists$} do


  # create the user
  user_uri = URI('http://127.0.0.1:9095/setup/create_initial_user')
  req = Net::HTTP::Post.new(user_uri, "Content-Type" => "application/json")
  req.body = user.to_json
  res = Net::HTTP.start(user_uri.hostname, user_uri.port) do |http|
    http.request(req)
  end
  expect(res.code_type).to eql(Net::HTTPCreated)

  # FIXME add steps for setting public key once complete
  step "I have logged in"
end

Given %r{^my private key is valid$} do
  step "I have an SSH key"
  @ssh_key = valid_key
end

Given %r{^my private key is invalid$} do
  @ssh_key = invalid_key
end

When %r{^I initiate a SSH connection to the server$} do
  @response = connect_via_ssh(key: @ssh_key)
end

Then %r{^I should see the response "(.*?)"$} do |text|
  expect(@response).to include(text)
end

Then %r{^I should receive the error message "(.*?)"$} do |message|
  step %{I should see the response "#{message}"}
end

Then %r{^I should receive an AuthenticationFailed error$} do
  expect(@response).to be_a(Net::SSH::AuthenticationFailed)
end

When %r{^I try to run one of the following commands:$} do |table|
  @commands = table.transpose.raw.flatten
end

Then %r{^I should receive an error$} do
  @commands.each do |command|
    connect_via_ssh(cmd: command).tap do |response|
      expect(response).to include("Only Git operations are permitted")
    end
  end
end

When %r{^I try to clone the repository "(.*?)"$} do |name|

  # First make sure the target dir doesn't exist
  @clone_dir = "content"
  @clone_location = "../tmp/ssh/#{@clone_dir}"
  FileUtils.rm_rf(@clone_location)
  expect(Dir.exists?(@clone_location)).to be false

  Dir.chdir("../tmp/ssh") do
    cert_path = "../../data/certificates/valid"
    config_path = "../../data/ssh/config"
    command = "git clone git@127.0.0.1:#{name}"
    full_command = %{GIT_SSH_COMMAND='ssh -F #{config_path} -i #{cert_path}' #{command} #{@clone_dir}}

    # FXIME, this outputs text, should really surpress within Cucumber (or send it straight
    # to the log)
    @blah, @output, @status = Open3.capture3(full_command)
  end

end

Then %r{^I should see output detailing my clone operation$} do
  expect(@output).to include("Cloning into '#{@clone_dir}'")
end

Then %r{^the directory should be present in my working directory$} do
  expect(Dir.exists?(@clone_location)).to be true
end

Given %r{^I try to establish a connection with user "(.*?)"$} do |username|
  @response = connect_via_ssh(user: username)
end

Then %r{^I should see output with an error$} do
  expect(@output).to include("fatal: protocol error")
end

def connect_via_ssh(host: "127.0.0.1", port: 2223, key: valid_key, cmd: "", user: "git")

  response = StringIO.new

  Net::SSH.start(
    host, user,
    port: port,
    host_key: "ssh-rsa",
    #encryption: "aes128-cbc",
    #compression: "zlib",
    keys: [key]
  ) do |session|

    session.open_channel do |channel|
      channel.on_data do |ch, data|
        response << "#{data}"
      end
      channel.exec(cmd)
    end

    session.loop
  end

  return response.string

rescue Net::SSH::AuthenticationFailed => e
  return e
end


def valid_key
  "../data/certificates/valid".tap do |path|
    File.chmod(0600, path)
  end
end

def invalid_key
  "../data/certificates/missing"
end
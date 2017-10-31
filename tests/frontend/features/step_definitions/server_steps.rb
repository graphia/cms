Given %r{^the CMS is running with the "(.*?)" config$} do |opt|

  config_mapping = {
    "default" => "../../config/cucumber.yml"
  }

  start_server(config_mapping[opt])
end

def start_server(server_config_path)

  # kill existing pid first
  if FileTest.exist?(PID_PATH)
    begin
      @pid = Pathname.new(PID_PATH).read.to_i

      kill(@pid) if Process.getpgid(@pid) && @pid > 0
    rescue Errno::ESRCH
      # already dead 😵
    end
  end

  @pid = Process.spawn(
    [
      "../../graphia-cms",
      "-config #{server_config_path}",
      "-log-to-file true"
    ].join(" ")
  )

  # Create a download directory if one doesn't already exist
  # if it does exist, make sure it's empty
  if Dir.exist?(DOWNLOAD_DIR)
    FileUtils.rm_rf(Dir.glob("#{DOWNLOAD_DIR}/*"))
  else
    FileUtils.mkdir(DOWNLOAD_DIR)
  end

    #command = "../../graphia-cms -config=../../config/cucumber.yml -log-to-file=true "

    #Open3.popen3(command) do |stdin, stdout, stderr, wait_thr|
      # FIXME Negroni's output is still appearing, work out how to suppress it
    #end

  Pathname.new(PID_PATH).write(@pid)

  # wait for the server to be running before contining
  1.upto(WEB_SERVER_START_ATTEMPTS) do |inc|
    if inc == WEB_SERVER_START_ATTEMPTS
      fail "Timed out waiting for web server to start"
    else
      break if attempt_webserver_call
      sleep 0.1
    end
  end

end

def	attempt_webserver_call
  Net::HTTP.get('localhost', '/', '9095')
  true
rescue Errno::ECONNREFUSED
  false
end
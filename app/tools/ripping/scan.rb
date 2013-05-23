dvd_drive = 'F:'
scan_command = "HandBrakeCLI.exe --scan --title 0 -i #{dvd_drive}\\ 2>&1"
output_dir = "C:\\wan\\rips"
minimum_duration = 1200

def format_ep_num(num)
  num.to_s.rjust(2).sub(" ", "0")
end

def slugify(s)
  words = s.split(/[ _-]/).collect {|x| x.downcase}
  words.collect! {|word| word = word[0].upcase + word[1..-1]}
  return words.join(" ")
end

def parse_arg(args, patterns, expect_value=false, default_value=nil)
  if (!patterns.kind_of?(Array))
    patterns = [patterns]
  end

  args.each_index do |i|
    patterns.each do |pattern|
      if (args[i] == pattern)
        if (expect_value)
          return args[i+1] if (args.length > i+1)
        else
          return true
        end
      end
    end     
  end
  
  return default_value
end

def usage
  "ruby scan.rb [--start-ep-num|-ep EPISODE_NUMBER] [--split-single|-ss EPISODE_COUNT] \\
  [--output-dir|--output|-o PATH] [-i|--input|--dvd-drive DRIVE_LETTER:] \\
  [--subtitle-track|-s TRACKNUM] [--sub-audio|--subtitle-audio AUDIO_TRACK_NUM] \\
  [--title-override TITLE_ROOT] [--scan]"
end

## Step 0: Parse the commandline arguments
start_episode_num = parse_arg(ARGV, ['--start-ep-num', '-ep'], true, 1).to_i
split_single_count = parse_arg(ARGV, ['--split-single', '-ss'], true).to_i
output_dir = parse_arg(ARGV, ['--output-dir', '-o', '--output'], true, "C:\\wan\\rips")
dvd_drive = parse_arg(ARGV, ['-i', '--input', '--dvd-drive'], true, 'F:')
force_subtitle_track = parse_arg(ARGV, ['--subtitle-track', '-s'], true)
force_sub_audio = parse_arg(ARGV, ['--sub-audio', '--subtitle-audio'], true)
scan_only = parse_arg(ARGV, ['--scan'], false, false)
title_override = parse_arg(ARGV, ['--title-override'], true, nil)
ep_num_join_char = parse_arg(ARGV, ['--join-string', '-js'], true, 'E')

subtitle_track = force_subtitle_track
sub_audio = force_sub_audio
  
if ($DEBUG)
  puts "Config:"
  puts "-Split Single: #{split_single_count}"
  puts "-Starting Episode Number: #{start_episode_num}"
  puts "-Output Dir: #{output_dir}"
  puts "-DVD Drive: #{dvd_drive}"
  puts "-Title Root Override: #{title_override}"
  puts "-Subtitle Track: #{force_subtitle_track}"
  puts "-Sub Audio: #{force_sub_audio}"
  puts "-Scan Only? #{scan_only}"
end

if (ARGV.include?('--help'))
  puts usage
  exit 0
end

class Title
    attr_accessor :num, :chapters, :duration, :size, :audio, :subtitles

    def initialize(num='1', chapters=[], duration='1:00', size='720x480',
                    audio= {}, subtitles = {})
        @num = num
        @chapters = chapters
        @duration = duration
        @size = size
        @audio = audio
        @subtitles = subtitles
    end

    def duration_seconds
        tokens = @duration.split(/:/).collect {|x| x.to_i}
        seconds = tokens[0] * 3600 + tokens[1] * 60 + tokens[2]
        return seconds
    end

    def to_s
        s = "Title ##{@num}, #{self.duration_seconds.to_s} seconds (#{@duration})"
        s += "\n Audio Tracks: " + @audio.collect {|track, language| "#{track} => #{language}"}.join(", ")
        s += "\n Subtitle Tracks: " + @subtitles.collect {|track, language| "#{track} => #{language}"}.join(", ")
        s += "\n Recommend using subtitle mode? #{self.subtitle_eligible.to_s}"
        s += "\n Chapters: #{@chapters.length}"
    end
    
    def subtitle_eligible
      jp_audio = false
      @audio.each_pair do |track, language|
        if (language == "Japanese")
          jp_audio = track
          break
        end
      end
      
      en_subtitle = false
      @subtitles.each_pair do |track, language|
        if (language == "English")
          en_subtitle = track
          break
        end
      end
      
      if (jp_audio == false || en_subtitle == false)
        return false
      end
      
      return [jp_audio, en_subtitle]
    end

    def Title.parse(scan = [])
        t = Title.new

        if (scan.length < 5)
            puts "Not enough data (#{scan.length} lines) - this can't be a successful scan. Aborting..."
            exit
        end

        t.num= scan[0].strip.split(/[ :]/)[-1].to_i
        scan.each_index do |line_i|
            line = scan[line_i].strip
            if (line =~ /\+ duration/) 
                t.duration = line.strip.split(/ /)[-1]
            elsif (line =~ /size/)
                t.size = line.split(/[ ,]/)[2]
            elsif (line =~ /audio tracks/)
              # Parse out all of the following lines that have language data
              (line_i + 1).upto(scan.length-1) do |j|
                # Quit looking when the first non-language line is found (this means the next section is starting)
                if (scan[j] =~ /subtitle tracks/)
                  break
                end
                  
                # Parse out the track number and language for each language track
                tokens = scan[j].strip.split(/[ ,]+/)
                t.audio[tokens[1]] = tokens[2]
              end
            elsif (line =~ /subtitle tracks/)
                # Parse out all of the following lines that have language data
              (line_i + 1).upto(scan.length-1) do |j|
                # Quit looking when the first non-language line is found (this means the next section is starting)
                if (scan[j] =~ /^\+title /)
                  break
                end
                  
                # Parse out the track number and language for each language track
                tokens = scan[j].strip.split(/[ ,]+/)
                t.subtitles[tokens[1]] = tokens[2]
              end
            elsif (line =~ /\+ \d+: cells /)
              tokens = line.split(/[ ,]/).collect {|x| (x.length == 0) ? nil : x}.compact
              chapter_number = tokens[1][0...-1].to_i
              duration = tokens[-1]
              t.chapters[chapter_number] = duration
            end
        end
        return t
    end
end

##############
#### MAIN ####
##############

puts "Split-single with #{split_single_count} chapters" if($DEBUG)

## Step I: Scan the Disk's Titles
scan_output = `#{scan_command}`.split(/\n/)
useful_output = []
title_lines = []
useful_line_count = 0

current_title = 0
accumulated_lines = []
title_blocks = {}

scan_output.each do |line|
  if (line =~ /^([\s]*)\+ /)
#    puts line
    if (line =~ /^\+ title \d+:/)
      if (current_title > 0)
        title_blocks[current_title] = accumulated_lines
		  end
		  current_title = line.strip.split(/[ :]/)[2].to_i
		  accumulated_lines = []
    end
    accumulated_lines << line
    useful_line_count += 1
  end
end

## Step II: Parse the scan data 
titles = []
title_blocks.each_pair do |title, scan_data|
	t = Title.parse(scan_data)
	titles << t
end

## Step III: Read the disk's volume label to use for the default filenames
out = `dir #{dvd_drive}`.split(/\n/).collect {|line| line.strip}
volume_label = slugify(out[0].split(/ /)[5..-1].join(' '))

puts "Using DVD Volume label for default output names: #{volume_label}"


## Step IV: Rip each Title that's over 20 minutes long

# Do some preprocessing.
#  IV a) First reduce the titles to just the ones over the minimum duration
titles_to_rip = []
titles.each do |title|
  if (title.duration_seconds >= minimum_duration)
    titles_to_rip << title
  end
end
# Save the filtered title list
titles = titles_to_rip
titles_to_rip = nil

puts "Using #{titles.length} titles over #{minimum_duration} seconds long"
titles.each do |title|
  puts title.to_s
end

# IV b) Next decide whether to use the default language, or Japanese+subtitles
#    subbing is enabled when all of the eligible tracks have both a Japanese audio track
#    as well as an English subtitle track.
use_subtitles = []
titles.each do |title|
  use_subtitles << title.subtitle_eligible
end

puts "Subtitling analysis:"
puts use_subtitles.collect {|sub_data| (sub_data.kind_of?(Array)) ? "[ #{sub_data[0]}, #{sub_data[1]}]" : sub_data.to_s}.join(", ")

if (!use_subtitles.index(false).nil?)
  use_subtitles = false
else
  sub_audio = (force_sub_audio.nil?) ? use_subtitles[0][0] : force_sub_audio 
  subtitle_track = (force_subtitle_track.nil?) ? use_subtitles[0][1] : force_subtitle_track
end


# IV c) Actually rip each title
# determine if there is a single hour+ duration track
big_track_index = nil
title_root = (title_override.nil?) ? volume_label : title_override
titles.each_index do |i|
  puts "Looking for a big track: [#{i}] => #{titles[i].duration}" if ($DEBUG)
  if (titles[i].duration_seconds > 3600)
    big_track_index = i
    break
  end
end

if (!big_track_index.nil?)
  puts "Found big track at index #{big_track_index}"
end

if (!(split_single_count.nil? || split_single_count <= 0) && !big_track_index.nil?)
  
  # Split the single-title into multiple episodes as best we can
  big_title = titles[big_track_index]
  chapters_per_episode = big_title.chapters.length / split_single_count
  extra_chapters = big_title.chapters.length % split_single_count
  
  puts "Ripping in split-single mode on track #{big_track_index+1} with #{big_title.chapters.length} chapters, total duration=#{big_title.duration}"
  puts "-- Chapters/episode: #{chapters_per_episode}, Extra Chapters: #{extra_chapters}" if ($DEBUG)
  
  rip_cmd_base = "HandBrakeCLI.exe -i F:\\ --title #{big_title.num}"
  rip_cmd_base += " --subtitle-burn --audio #{sub_audio} --subtitle #{subtitle_track}" if (use_subtitles)
  
  0.upto(split_single_count-1) do |i|
    start_chapter = (i * chapters_per_episode) + 1
    end_chapter = start_chapter + chapters_per_episode - 1
    chapter_arg = " -c #{start_chapter}-#{end_chapter}"
    
    rip_cmd = String.new(rip_cmd_base)
    rip_cmd += chapter_arg
    
    outpath = "#{output_dir}\\#{title_root}#{ep_num_join_char}#{format_ep_num(start_episode_num)}.m4v"
    start_episode_num += 1
    rip_cmd += " -o \"#{outpath}\""
    
    puts " > #{rip_cmd}"
    `#{rip_cmd} 2>autorip.error.log >> autorip.log` unless($DEBUG || scan_only)
  end
      
    
else
  # Rip each title individually
  puts "Ripping in single-track mode on #{titles.length} titles"
  titles.each do |title|
    output_path = "#{output_dir}\\#{title_root}#{ep_num_join_char}#{format_ep_num(start_episode_num)}.m4v"
    start_episode_num += 1
    rip_cmd = "HandBrakeCLI.exe -i F:\\ --title #{title.num} -o \"#{output_path}\""
    if (use_subtitles)
      rip_cmd += " --subtitle-burn --audio #{sub_audio} --subtitle #{subtitle_track}"
    end
    puts "Ripping Title ##{title.num} to file #{output_path}"
		
    puts " > #{rip_cmd}"
    `#{rip_cmd} 2>autorip.error.log >> autorip.log` unless($DEBUG || scan_only)
  end
end

## Step V: Eject the DVD
`C:\\wan\\bin\\ejectcd.exe` unless($DEBUG || scan_only)



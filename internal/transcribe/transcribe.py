import whisper
import sys
import urllib.request
import os
import uuid
import tempfile

import whisper
import sys
import requests
import os
import uuid
import tempfile

def download_audio(url):
    temp_dir = tempfile.gettempdir()
    filename = os.path.join(temp_dir, f"{uuid.uuid4()}.mp3")
    
    try:
        print(f"Downloading from: {url}", file=sys.stderr)
        
        # Dùng requests với verify=False
        response = requests.get(url, verify=False)
        response.raise_for_status()
        
        with open(filename, 'wb') as f:
            f.write(response.content)
            
        print(f"Downloaded to: {filename}", file=sys.stderr)
        return filename
    except Exception as e:
        print(f"Error downloading: {e}", file=sys.stderr)
        raise

    
if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Audio URL required", file=sys.stderr)
        sys.exit(1)

    url = sys.argv[1]
    
    try:
        temp_path = download_audio(url)
        
        print("Loading Whisper model...", file=sys.stderr)
        model = whisper.load_model("small")
        
        print("Transcribing...", file=sys.stderr)
        result = model.transcribe(temp_path)
        
        # Output chính (stdout)
        print(result["text"])
        
        # Cleanup
        os.remove(temp_path)
        
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)
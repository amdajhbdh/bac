# BAC Unified - Storage Module
# Garage S3-compatible storage for media files

# Storage configuration
export def get-config [] {
    let config_path = "config/storage.yaml"
    
    if not ($config_path | path exists) {
        return {
            endpoint: "http://localhost:3900"
            region: "us-east-1"
            access_key: "garage_test_key"
            secret_key: "garage_test_secret"
            bucket: "bac-unified"
            data_dir: "./data/storage"
        }
    }
    
    open $config_path | get storage
}

# Check if Garage is available
export def is-available [] {
    let config = (get-config)
    
    try {
        let response = (http head $config.endpoint -e)
        $response.status == 200
    } catch {
        false
    }
}

# Upload file to Garage
export def upload [file: string, object-name: string] {
    let config = (get-config)
    let file_path = (ls -la $file | get name.0)
    
    if not ($file_path | path exists) {
        print $"File not found: ($file)"
        return null
    }
    
    # Using curl to upload to Garage S3 API
    let url = $"($config.endpoint)/($config.bucket)/($object-name)"
    
    try {
        let content_type = (match ($file | path parse | get extension) {
            "jpg" => "image/jpeg",
            "jpeg" => "image/jpeg",
            "png" => "image/png",
            "gif" => "image/gif",
            "pdf" => "application/pdf",
            "mp4" => "video/mp4",
            "mp3" => "audio/mpeg",
            _ => "application/octet-stream"
        })
        
        let status = (http put $url -d [($file | path read --raw)] --content-type $content-type -H [ $"Authorization: AWS4-HMAC-SHA256 Credential=($config.access_key)" ] -s | complete | get exit_code)
        
        if $status == 0 {
            print $"Uploaded: ($object-name)"
            return $url
        }
    } catch {
        print $"Upload failed: ($in)"
    }
    
    # Local fallback
    let dest = ($config.data_dir | path join $object-name)
    mkdir ($dest | path dirname) -p
    cp $file $dest
    print $"Saved locally: ($dest)"
    return $dest
}

# Download file from Garage
export def download [object-name: string, dest: string] {
    let config = (get-config)
    let url = $"($config.endpoint)/($config.bucket)/($object-name)"
    
    try {
        http get $url -o $dest
        print $"Downloaded: ($dest)"
        return $dest
    } catch {
        # Local fallback
        let local_path = ($config.data_dir | path join $object-name)
        if ($local_path | path exists) {
            cp $local_path $dest
            print $"Downloaded from local: ($dest)"
            return $dest
        }
        print $"File not found: ($object-name)"
        return null
    }
}

# List objects in bucket
export def list [prefix: string = ""] {
    let config = (get-config)
    
    # Try Garage API
    let url = $"($config.endpoint)/($config.bucket)?prefix=($prefix)"
    
    try {
        let response = (http get $url -e | get body)
        print $response
    } catch {
        # Local fallback
        let local_dir = ($config.data_dir | path join $prefix)
        if ($local_dir | path exists) {
            ls -la $local_dir
        } else {
            print "No files found"
        }
    }
}

# Delete object
export def delete [object-name: string] {
    let config = (get-config)
    let url = $"($config.endpoint)/($config.bucket)/($object-name)"
    
    try {
        http delete $url -e
        print $"Deleted: ($object-name)"
    } catch {
        # Local fallback
        let local_path = ($config.data_dir | path join $object-name)
        if ($local_path | path exists) {
            rm $local_path
            print $"Deleted locally: ($object-name)"
        }
    }
}

# Get object URL
export def get-url [object-name: string] {
    let config = (get-config)
    $"($config.endpoint)/($config.bucket)/($object-name)"
}

# Storage status
export def status [] {
    print "=== Storage Status ==="
    
    let config = (get-config)
    print $"Endpoint: ($config.endpoint)"
    print $"Bucket: ($config.bucket)"
    print $"Data Dir: ($config.data_dir)"
    print ""
    
    if (is-available) {
        print "Garage: Available"
    } else {
        print "Garage: Unavailable (using local storage)"
    }
    
    # Count local files
    if ($config.data_dir | path exists) {
        let count = (ls -la $config.data_dir | length)
        print $"Local files: ($count)"
    }
}

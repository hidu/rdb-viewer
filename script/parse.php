<?php
/***************************************************************************
 * 
 * Copyright (c) 2016 github.com/hidu, Inc. All Rights Reserved
 * 
 **************************************************************************/
 
 
 
/**
 * @file parse.php
 * @author hidu
 * @date 2016/05/13 22:00:27
 * @brief 
 *
 *rdb_viewer -json -val hd.rdb|php script/parse.php 
 *
 **/

while(!feof(STDIN)){
    $line=fgets(STDIN);
	$obj=json_decode($line,true);
	if(is_array($obj)){
	    foreach ($obj as $k=>$v){
	        if(strpos($k,"_b")){
	            $obj[substr($k, 0,strlen($k)-2)]=base64_decode($v);
	            unset($obj[$k]);
	        }
	    }
	}
	//your code
	print_r($obj);
}

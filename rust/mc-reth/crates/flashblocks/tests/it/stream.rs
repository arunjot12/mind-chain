use futures_util::stream::StreamExt;
use reth_mind_chain_flashblocks::WsFlashBlockStream;

#[tokio::test]
#[ignore = "depends on external service availability"]
async fn test_online_streaming_flashblocks_from_remote_source_is_successful() {
    let items = 3;
    let ws_url = "wss://sepolia.flashblocks.base.org/ws".parse().unwrap();
    let stream = WsFlashBlockStream::new(ws_url);

    let blocks: Vec<_> = stream.take(items).collect().await;

    for block in blocks {
        assert!(block.is_ok());
    }
}

package go;

// http://www.javamex.com/java_equivalents/unsigned_arithmetic.shtml

public final class Unsigned{
	
	public static final int div32(int x, int y){
  		return (int) ((long)(x & 0xFFFFFFFFL) / (long)(y & 0xFFFFFFFFL));
	}

	public static final short div16(short x, short y){
  		return (short) ((int)(x & 0xFFFF) / (int)(y & 0xFFFF));
	}

	public static final byte div8(byte x, byte y){
  		return (byte) ((int)(x & 0xFF) / (int)(y & 0xFF));
	}
	
}

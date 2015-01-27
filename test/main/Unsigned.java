package main;

// http://www.javamex.com/java_equivalents/unsigned_arithmetic.shtml

public final class Unsigned{
	
	public static final int quo32(int x, int y){
  		return (int) ((long)(x & 0xFFFFFFFFL) / (long)(y & 0xFFFFFFFFL));
	}

	public static final short quo16(short x, short y){
  		return (short) ((int)(x & 0xFFFF) / (int)(y & 0xFFFF));
	}

	public static final byte quo8(byte x, byte y){
  		return (byte) ((int)(x & 0xFF) / (int)(y & 0xFF));
	}


	public static final boolean lss32(int x, int y){
  		return ((long)(x & 0xFFFFFFFFL) < (long)(y & 0xFFFFFFFFL));
	}

	public static final boolean lss16(short x, short y){
  		return ((int)(x & 0xFFFF) < (int)(y & 0xFFFF));
	}

	public static final boolean lss8(byte x, byte y){
  		return ((int)(x & 0xFF) < (int)(y & 0xFF));
	}
	

	public static final boolean leq32(int x, int y){
  		return ((long)(x & 0xFFFFFFFFL) <= (long)(y & 0xFFFFFFFFL));
	}

	public static final boolean leq16(short x, short y){
  		return ((int)(x & 0xFFFF) <= (int)(y & 0xFFFF));
	}

	public static final boolean leq8(byte x, byte y){
  		return ((int)(x & 0xFF) <= (int)(y & 0xFF));
	}
	

	public static final boolean gtr32(int x, int y){
  		return ((long)(x & 0xFFFFFFFFL) > (long)(y & 0xFFFFFFFFL));
	}

	public static final boolean gtr16(short x, short y){
  		return ((int)(x & 0xFFFF) > (int)(y & 0xFFFF));
	}

	public static final boolean gtr8(byte x, byte y){
  		return ((int)(x & 0xFF) > (int)(y & 0xFF));
	}
	

	public static final boolean geq32(int x, int y){
  		return ((long)(x & 0xFFFFFFFFL) >= (long)(y & 0xFFFFFFFFL));
	}

	public static final boolean geq16(short x, short y){
  		return ((int)(x & 0xFFFF) >= (int)(y & 0xFFFF));
	}

	public static final boolean geq8(byte x, byte y){
  		return ((int)(x & 0xFF) >= (int)(y & 0xFF));
	}
}
